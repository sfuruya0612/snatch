NAME := snatch
COMMIT_HASH := $(shell git rev-parse --short HEAD)
LDFLAGS := -X 'main.commit=${COMMIT_HASH}'

.PHONY: all
all: install

.PHONY: init
init:
	asdf install
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

.PHONY: prepare
prepare: fmt vet tidy

.PHONY: fmt
fmt:
	go fmt ./...

.PHONY: vet
vet:
	go vet ./...

.PHONY: tidy
tidy:
	go mod tidy -v -go=1.21

.PHONY: update
update:
	go get -u ./...

.PHONY: golangci-lint
golangci-lint:
	golangci-lint run ./...

.PHONY: test
test: prepare
	go test -v -race --cover ./...

.PHONY: install
install:
	go install -ldflags "${LDFLAGS}"
