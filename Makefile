DATE := $(shell TZ=Asia/Tokyo date +%Y%m%d-%H:%M:%S)
HASH := $(shell git rev-parse --short HEAD)
GOVERSION := $(shell go version)
LDFLAGS := -X 'main.date=${DATE}' -X 'main.hash=${HASH}' -X 'main.goversion=${GOVERSION}'

NAME := snatch
MODULE := github.com/sfuruya0612/${NAME}

AWS_PROFILE := default
REGION := ap-northeast-1

.PHONY: test build image

test:
	golangci-lint run
	go test -v --cover ./...

build: test
	-rm -rf build
	mkdir build

	GOOS=linux GOARGH=amd64 go build -ldflags "${LDFLAGS}" ${MODULE}
	zip build/${NAME}_linux_amd64.zip ${NAME}

	GOOS=darwin GOARGH=amd64 go build -ldflags "${LDFLAGS}" ${MODULE}
	zip build/${NAME}_darwin_amd64.zip ${NAME}

	@rm ${NAME}

image: build
	docker-compose build

install: test image
	-rm ${GOPATH}/bin/${NAME}
	go mod tidy
	go install -ldflags "${LDFLAGS}" ${MODULE}

clean:
	-rm ${GOPATH}/bin/${NAME}
	-rm -rf build
	-docker rmi --force ${NAME}_cli

# Testing
pip_install:
	pushd test ; pip install -r requirements.txt; popd

create_stack: pip_install
	python test/create_stack.py -a ${NAME} -p ${AWS_PROFILE} -r ${REGION} &

delete_stack: pip_install
	python test/delete_stack.py -a ${NAME} -p ${AWS_PROFILE} -r ${REGION} &
