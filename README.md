# snatch

[![Go Report Card](https://goreportcard.com/badge/github.com/sfuruya0612/snatch)](https://goreportcard.com/report/github.com/sfuruya0612/snatch)

## Description

Cli command to get and display Amazon Web Services resources.
This tool allows you to retrieve the content you need without logging into the management console.

The concept is that you can continue to work on checking resources without having to leave the terminal software.

## Installation

### go get

Install with `go get` command.

```sh
go get github.com/sfuruya0612/snatch
```

### Source build

Clone the repository and build it.

You need to install [mise](https://mise.jdx.dev/) to build.

```sh
git clone https://github.com/sfuruya0612/snatch.git
cd snatch
mise install
mise run install
```

## Shell completion

Enable auto-completion on tabs.
It is set to match $SHELL.

## Documentation

## License

[MIT License](./LICENSE)
