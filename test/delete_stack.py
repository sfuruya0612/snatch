# -*- coding: utf-8 -*-

import os
import sys
import logging
from boto3.session import Session
from botocore.exceptions import ClientError
from argparse import ArgumentParser

TEMPLATE = "test-ec2"

logger = logging.getLogger()
formatter = '%(levelname)s : %(asctime)s : %(message)s'
logging.basicConfig(level=logging.INFO, format=formatter)

class DeleteStack:

    # Option parser
    def get_option(self):
        usage = "python " + sys.argv[0] + " [-h | --help] [-a | --app <APP_NAME>] [-p | --profile <AWS_PROFILE>] [-r | --region <AWS_REGION>]"
        argparser = ArgumentParser(usage=usage)
        argparser.add_argument("-a", "--app", type=str,
                               default="snatch",
                               help="Target app name.")
        argparser.add_argument("-p", "--profile", type=str,
                               default="default",
                               help="~/.aws/config.")
        argparser.add_argument("-r", "--region", type=str,
                               default="ap-northeast-1",
                               help="AWs regions. e.g. ap-northeast-1, us-east-1, ...")
        return argparser.parse_args()

    def delete_stack(self, app_name, session):
        cfn = session.client("cloudformation")
        waiter = cfn.get_waiter("stack_delete_complete")

        stack_name = app_name + "-" + TEMPLATE

        logger.info("Target stack: %s", stack_name)
        try:
            cfn.delete_stack(
                StackName = stack_name,
            )

            waiter.wait(
                StackName = stack_name,
            )

            logger.info("%s delete complete.", stack_name)

        except ClientError as e:
            logger.warning("%s", e.response["Error"]["Message"])
            sys.exit(1)

    def delete_keypair(self, app_name, session):
        logger.info("Delete %s KeyPair.", app_name)

        ec2 = session.client("ec2")

        try:
            ec2.delete_key_pair(
                KeyName= app_name,
            )

            os.remove(os.environ["HOME"] + "/.ssh/" + app_name + ".pem")

            logger.info("%s KeyPair delete complete.", app_name)

        except ClientError as e:
            logger.warning("%s", e.response["Error"]["Message"])
            sys.exit(1)

    @staticmethod
    def main():
        logger.info("Start delete stacks.")

        self = DeleteStack()

        options = self.get_option()
        app_name = options.app
        profile = options.profile
        region = options.region

        session = Session(region_name=region, profile_name=profile)

        self.delete_stack(app_name, session)
        self.delete_keypair(app_name, session)

        logger.info("Finish delete stacks.\n")

if __name__ == '__main__':
    DeleteStack.main()
