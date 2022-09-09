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
	go mod tidy -v -go=1.19

.PHONY: golangci-lint
golangci-lint:
	golangci-lint run ./...

.PHONY: test
test: prepare
	go test -v -race --cover ./...

.PHONY: build
build: test
	-mkdir build

	GOOS=linux GOARGH=amd64 go build -ldflags "${LDFLAGS}"
	zip build/${NAME}_linux_amd64.zip ${NAME}

	GOOS=linux GOARGH=arm64 go build -ldflags "${LDFLAGS}"
	zip build/${NAME}_linux_arm64.zip ${NAME}

	GOOS=darwin GOARGH=amd64 go build -ldflags "${LDFLAGS}"
	zip build/${NAME}_darwin_amd64.zip ${NAME}

	GOOS=darwin GOARGH=arm64 go build -ldflags "${LDFLAGS}"
	zip build/${NAME}_darwin_arm64.zip ${NAME}

	@rm ${NAME}

.PHONY: install
install:
	-rm ${GOPATH}/bin/${NAME}
	go install -ldflags "${LDFLAGS}"

.PHONY: clean
clean:
	-rm ${GOPATH}/bin/${NAME}
	-rm -rf build
