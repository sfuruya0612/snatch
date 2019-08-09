# snatch
AWSリソース情報を取得するGolang製ツール(該当AWS CLI options：describe, event, ...)

## Help
``` sh
$ snatch -h
NAME:
   snatch - Show AWS resources cli command.

USAGE:
   snatch [global options] command [command options] [arguments...]

VERSION:
   YYYYMMDD-hh:mm:ss xxxxyyyy (go version go1.12.5 darwin/amd64)

COMMANDS:
     ec2      Show EC2 resources. (default: Describe EC2 instances)
     rds      Show RDS resources. (default: Describe RDS instances)
     ec       Show ElastiCache resources. (default: Describe Cache Clusters)
#    ~~~
#    Add a AWS services.
#    ~~~
     help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --profile value, -p value  Choose AWS credential. (default: "default")
   --region value, -r value   Select Region. (default: "ap-northeast-1")
   --help, -h                 show help
   --version, -v              print the version

```

## Install
``` sh
make install
```

## Docker run
``` sh
make image
docker-compose run cli snatch -p <value> <command>
```

## Uninstall
``` sh
make clean
```
