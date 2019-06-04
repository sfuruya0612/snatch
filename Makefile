DATE := $(shell TZ=Asia/Tokyo date +%Y%m%d-%H:%M:%S)
HASH := $(shell git rev-parse --short HEAD)
GOVERSION := $(shell go version)
LDFLAGS := -X 'main.date=${DATE}' -X 'main.hash=${HASH}' -X 'main.goversion=${GOVERSION}'

APP := snatch
MODULE := github.com/sfuruya0612/${APP}/cmd/${APP}

install:
	-rm ${GOPATH}/bin/${APP}
	go mod tidy
	go install -ldflags "${LDFLAGS}" ${MODULE}

delete:
	-rm ${GOPATH}/bin/${APP}

build:
	rm -rf build
	mkdir build

	GOOS=linux GOARGH=amd64 go build -ldflags "${LDFLAGS}" ${MODULE}
	zip build/${APP}_linux_amd64.zip ${APP}

	GOOS=darwin GOARGH=amd64 go build -ldflags "${LDFLAGS}" ${MODULE}
	zip build/${APP}_darwin_amd64.zip ${APP}

