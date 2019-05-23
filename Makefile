DATE := $(shell TZ=Asia/Tokyo date +%Y%m%dT%H%M%S+0900)
HASH := $(shell git rev-parse HEAD)
GOVERSION := $(shell go version)
LDFLAGS := -X 'main.date=${DATE}' -X 'main.hash=${HASH}' -X 'main.goversion=${GOVERSION}'

APP := snatch
MODULE := github.com/ShoichiFuruya/${APP}/cmd/${APP}
ROOT := ${GOPATH}/src/${MODULE}

install:
	-rm ${GOPATH}/bin/${APP}
	go mod tidy
	go install -ldflags "${LDFLAGS}" ${MODULE}

