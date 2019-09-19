# snatch
AWSリソース情報を取得するGolang製ツール  

## Getting started
* Go 1.11.x or later.
* You need to set your $GOPATH and have $GOPATH/bin in your path.

### Install
``` sh
go get github.com/sfuruya0612/snatch/cmd/snatch
```

or git clone (default install to darwin)
``` sh
git clone https://github.com/sfuruya0612/snatch
make install
snatch -h
```

### Setting the tab completion
``` sh
printf '\n%s\n%s\n%s\n' '# for snatch autocomplete' "test -f ~/.snatch_$(basename $SHELL)_autocomplete || curl -LRsS https://raw.githubusercontent.com/urfave/cli/master/autocomplete/$(basename $SHELL)_autocomplete -o ~/.snatch_$(basename $SHELL)_autocomplete" "PROG=snatch source ~/.snatch_$(basename $SHELL)_autocomplete" >> "${HOME}/.$(basename $SHELL)rc"
```

## Usage

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

## Testing

### Docker run(Testing linux ver)
``` sh
make image
docker-compose run cli snatch -p <value> <command>
```
