DATE := $(shell TZ=Asia/Tokyo date +%Y%m%d-%H:%M:%S)
HASH := $(shell git rev-parse --short HEAD)
GOVERSION := $(shell go version)
LDFLAGS := -X 'main.date=${DATE}' -X 'main.hash=${HASH}' -X 'main.goversion=${GOVERSION}'

NAME := snatch
MODULE := github.com/sfuruya0612/${NAME}

AWS_PROFILE := default
REGION := ap-northeast-1

.PHONY: test build image

init:
	asdf install
	go get -u github.com/kisielk/errcheck
	go get -u honnef.co/go/tools/cmd/staticcheck

test: init
	go fmt ./...
	go vet ./...
	errcheck ./...
	staticcheck ./...
	go test -v -race --cover ./...

build: test
	-rm -rf build
	mkdir build

	go mod tidy

	GOOS=linux GOARGH=amd64 go build -ldflags "${LDFLAGS}" ${MODULE}
	zip build/${NAME}_linux_amd64.zip ${NAME}

	GOOS=darwin GOARGH=amd64 go build -ldflags "${LDFLAGS}" ${MODULE}
	zip build/${NAME}_darwin_amd64.zip ${NAME}

	@rm ${NAME}

image: build
	docker-compose build

install: test
	-rm ${GOPATH}/bin/${NAME}
	go mod tidy
	go install -ldflags "${LDFLAGS}" ${MODULE}

clean:
	-rm ${GOPATH}/bin/${NAME}
	-rm -rf build
	-docker rmi --force ${NAME}_cli

# Test
pip_install:
	pushd scripts ; pip install -r requirements.txt; popd

deploy_stack: pip_install
	python scripts/deploy_stack.py -a ${NAME} -p ${AWS_PROFILE} -r ${REGION} &

delete_stack: pip_install
	python scripts/delete_stack.py -a ${NAME} -p ${AWS_PROFILE} -r ${REGION} &
