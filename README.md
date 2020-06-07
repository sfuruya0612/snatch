# snatch

<a href="https://goreportcard.com/report/github.com/sfuruya0612/snatch"><img src="https://goreportcard.com/badge/github.com/sfuruya0612/snatch" alt="Go Report Card"/></a>
<a href="https://travis-ci.org/sfuruya0612/snatch"><img src="https://travis-ci.org/sfuruya0612/snatch.svg?branch=master"></a>

## Index

- [Description](#description)
- [Getting started](#getting-started)
- [Usage](#usage)
- [Trial](#trial)

## Description

Cli command to get and display Amazon Web Services resources.  
This tool allows you to retrieve the content you need without logging into the management console.  

The concept is that you can continue to work on checking resources without having to leave the terminal software.

## Getting started

### Supported versions / OS

| OS / Laungage | Versions                                           |
| :------------ | :------------------------------------------------- |
| OS            | MacOS </br> Linux                                  |
| Go version    | 1.11.x </br> 1.12.x </br> 1.13.x </br> more latest |

### Install

Install with `go get` command.

``` sh
go get github.com/sfuruya0612/snatch
```

Compiling from Source.

``` sh
git clone https://github.com/sfuruya0612/snatch.git ./ && cd snatch && make install
```

Installed as a Docker image.

``` sh
git clone https://github.com/sfuruya0612/snatch.git ./ && cd snatch && make image
```

### Required settings

You need to set up AWS credential.  
Use the `aws configure` command it as needed.

``` sh
# e.g.
aws configure --profile myapp
```

### Optional settings

Enable auto-completion on tabs.  
It is set to match $SHELL.

``` sh
printf '\n%s\n%s\n%s\n' '# for snatch autocomplete' "test -f ~/.snatch_$(basename $SHELL)_autocomplete || curl -LRsS https://raw.githubusercontent.com/urfave/cli/master/autocomplete/$(basename $SHELL)_autocomplete -o ~/.snatch_$(basename $SHELL)_autocomplete" "PROG=snatch source ~/.snatch_$(basename $SHELL)_autocomplete" >> "${HOME}/.$(basename $SHELL)rc"
```

## Usage

How to use each command and an example of its execution.  
Details of the command can be found in `snatch -h`.

### EC2

``` sh
# Returns list of EC2 Instances
# Search by specifying Tags
$ snatch ec2
$ snatch ec2 --tag Name:*prod*

# Get EC2 system log (Output /var/log/cloud-init-output.log)
$ snatch ec2 log --id <YOUR INSTANCE ID>

# Terminate Instance
# Interactive confirmation at execute
$ snatch ec2 terminate --id <YOUR INSTANCE ID>
```

### RDS

``` sh
# Returns list of RDS Instances
$ snatch rds

# Returns list of RDS clusters
$ snatch rds cluster

```

### Elasticache

``` sh
# Returns list of Elasticache Clusters
$ snatch elasticache
$ snatch ec

# Returns list of Elasticache Nodes
$ snatch elasticache node
$ snatch ec node
```

### S3

``` sh
# Returns list of S3 Buckets
$ snatch s3

# Returns list of S3 Objects
# If you don't specify a bucket name, you can choose from a list of buckets
$ snatch s3 object
$ snatch s3 object --bucket <YOUR BUCKET NAME>

# Display S3 Object
$ snatch s3 cat --bucket <YOUR BUCKET NAME> --key <YOUR OBJECT KEY>

# Download S3 Object
$ snatch s3 cat --bucket <YOUR BUCKET NAME> --key <YOUR OBJECT KEY>
```

### on Docker 

``` sh
docker-composee run cli snatch <COMMANDS>
```

## Trial

### Setting

**To run python scripts, you need python and pip**  

Launch EC2 in aws cloudformation.  
The cost is minimal due to the use of EC2 Spot Instance, but there are running costs associated with it.  

#### Create resources

Create cloudformation stack.

``` sh
# Run on the backend
make create_stack
```

##### Executed commands

This is a list of commands that can be executed in this trial.

- EC2 list

``` sh
$ snatch ec2
```

- Login to EC2 with SSM

``` sh
$ snatch ssm
```

- Execute a command to EC2 using SSM Run Command

``` sh
$ snatch ssm run -t Name:snatch-test-pub-ec2 <COMMANDS>

# e.g.
$ snatch ssm run -t Name:snatch-test-pub-ec2 "date"
```

#### Delete resources

Delete cloudformation stack.

``` sh
# Run on the backend
make delete_stack
```
