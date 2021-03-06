GROUP := github.com/pmoncadaisla
NAME := go-journey

DOCKER_IMAGE := ${NAME}
PKG := ${GROUP}/${NAME}

OUT_BIN := build/app


PKG_LIST := $(shell go list ${PKG}/... | grep -v /vendor/)

VERSION := $(shell git describe --always --long --dirty || date)
GO_FILES := $(shell find . -name '*.go' | grep -v /vendor/)

all: build

dockerbuild: build
	@docker build -t ${DOCKER_IMAGE} -f Dockerfile .

build: test
	mkdir -p build
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -i -v -o ${OUT_BIN} -ldflags="-s -w -X main.Version=${VERSION}" ${PKG}

build-darwin: test
	mkdir -p build
	GOOS=darwin GOARCH=amd64 CGO_ENABLED=0 go build -a -i -v -o ${OUT_BIN}.darwin -ldflags="-s -w -X main.Version=${VERSION}" ${PKG}


test: 
	@go test -short ${PKG_LIST}

vet:
	@go vet ${PKG_LIST}


clean:
	-@rm ${OUT_BIN}

.PHONY: run build vet lint build