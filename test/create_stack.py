# -*- coding: utf-8 -*-

import glob
import os
import re
import sys
import logging
from boto3.session import Session
from botocore.exceptions import ClientError
from argparse import ArgumentParser

TEMPLATES = [
    "/test/ec2.yml",
]

logger = logging.getLogger()
formatter = '%(levelname)s : %(asctime)s : %(message)s'
logging.basicConfig(level=logging.INFO, format=formatter)

class CreateStack:

    # Option parser.
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

    # Create EC2 keypair.
    # 秘密鍵は ~/.ssh/ 配下に書き出す(file permission: 0600)
    def create_keypair(self, app_name, session):
        logger.info("Create %s KeyPair.", app_name)

        ec2 = session.client("ec2")

        try:
            ec2.describe_key_pairs(
                KeyNames=[
                    app_name,
                ],
            )

            logger.info("%s KeyPair already exists.", app_name)

        except ClientError as e:
            if e.response["Error"]["Code"] == "InvalidKeyPair.NotFound":
                res = ec2.create_key_pair(
                    KeyName=app_name,
                )

                logger.info("%s KeyPair Created.", app_name)

                private_key = res["KeyMaterial"]
                pem_file = open(os.environ["HOME"] + "/.ssh/" + app_name + ".pem", "w")
                pem_file.write(private_key)
                pem_file.close

                os.chmod(os.environ["HOME"] + "/.ssh/" + app_name + ".pem", 0o600)
            else:
                logger.warning("%s", e.response["Error"]["Message"])
                sys.exit(1)

    # Provisiond stack
    def provisiond(self, app_name, profile, region):
        session = Session(profile_name=profile, region_name=region)
        cfn = session.client("cloudformation")

        create_waiter = cfn.get_waiter("stack_create_complete")

        self.create_keypair(app_name, session)

        for template in TEMPLATES:
            path = os.getcwd() + template
            body = open(path).read()
            stack_name = app_name + "-" + re.sub('\/(.*)\/(.*)\.yml', '\\1-\\2', template)

            self.valid_template(template, body, cfn)

            try:
                self.create_stack(app_name, body, stack_name, cfn, create_waiter)

            except ClientError as e:
                logger.warning("%s", e.response["Error"]["Message"])
                sys.exit(1)

    # Valid CFn template.
    def valid_template(self, template, body, cfn):
        logger.info("Validate check %s", template)

        try:
            cfn.validate_template(
                TemplateBody = body,
            )

            logger.info("%s is validation OK.", template)

        except ClientError as e:
            logger.warning("%s", e.response["Error"]["Message"])
            sys.exit(1)

    # Create CFn stacks.
    def create_stack(self, app_name, body, stack_name, cfn, create_waiter):
        logger.info("Create %s.", stack_name)

        input = {
            "StackName": stack_name,
            "TemplateBody": body,
            "Capabilities": [
                'CAPABILITY_NAMED_IAM',
            ],
            "Parameters": [
                {
                    "ParameterKey": "AppName",
                    "ParameterValue": app_name,
                },
            ],
        }

        try:
            cfn.create_stack(**input)

            create_waiter.wait(
                StackName = stack_name,
            )

            logger.info("Create %s Complete.", stack_name)

        except ClientError as e:
            logger.warning("%s", e.response["Error"]["Message"])
            sys.exit(1)

    @staticmethod
    def main():
        logger.info("Start provision stacks.")

        self = CreateStack()

        options = self.get_option()
        app_name = options.app
        profile = options.profile
        region = options.region

        self.provisiond(app_name, profile, region)

        logger.info("Finish provision stacks.")

if __name__ == '__main__':
    CreateStack.main()
