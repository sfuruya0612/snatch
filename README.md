# snatch

<a href="https://goreportcard.com/report/github.com/sfuruya0612/snatch"><img src="https://goreportcard.com/badge/github.com/sfuruya0612/snatch" alt="Go Report Card"/></a>

Cli command to get and display Amazon Web Services resources.  
This tool allows cli to retrieve the necessary contents without logging into the Management console.  
The concept is that you can continue working without leaving the black screen (Terminal software).  

## Getting started

### Required
* Go version 1.11.x or later.
* You need to set your $GOPATH and have $GOPATH/bin in your path.

### Install
``` sh
go get github.com/sfuruya0612/snatch
```

or git clone
``` sh
git clone https://github.com/sfuruya0612/snatch.git
make install
```

### Setting the tab completion
``` sh
printf '\n%s\n%s\n%s\n' '# for snatch autocomplete' "test -f ~/.snatch_$(basename $SHELL)_autocomplete || curl -LRsS https://raw.githubusercontent.com/urfave/cli/master/autocomplete/$(basename $SHELL)_autocomplete -o ~/.snatch_$(basename $SHELL)_autocomplete" "PROG=snatch source ~/.snatch_$(basename $SHELL)_autocomplete" >> "${HOME}/.$(basename $SHELL)rc"
```

## Usage

``` sh
$ snatch -h
NAME:
   snatch - Cli command to get and display Amazon Web Services resources.

USAGE:
   snatch [global options] command [command options] [arguments...]

VERSION:
   20191011-11:36:07 0ae46416 (Build by: go version go1.12.5 darwin/amd64)

COMMANDS:
     ec2      Get a list of EC2 resources. (API: DescribeInstances)
     rds      Get a list of RDS resources. (API: DescribeDbInstances)
     ec       Get a list of ElastiCache Cluster resources. (API: DescribeCacheClusters)
     elb      Get a list of ELB(Classic) resources. (API: DescribeLoadBalancers)
     elbv2    Get a list of ELB(Application & Network) resources. (API: DescribeLoadBalancers)
     route53  Get a list of Rotue53 Record resources. (API: ListHostedZones and ListResourceRecordSets)
     acm      Get a list of ACM resources. (API: ListCertificates and DescribeCertificate)
     s3       Get Objects in selected S3 Bucket at interactive prompt. (API: ListBuckets and ListObjects)
     ssm      Start a session on your instances by launching bash or shell terminal. (API: StartSession)
     help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --profile value, -p value  Specify the AWS profile listed in ~/.aws/config. (default: "default") [$AWS_PROFILE]
   --region value, -r value   Specify the AWS region. (default: "ap-northeast-1")
   --help, -h                 show help
   --version, -v              print the version
```

## Testing

### Docker run (Testing linux version)
``` sh
make image
docker-compose run cli snatch -p <value> <command>
```
