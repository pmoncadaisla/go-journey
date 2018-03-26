GROUP := github.com/pmoncadaisla
NAME := go-journey

DOCKER_IMAGE := ${NAME}
PKG := ${GROUP}/${NAME}

OUT_BIN := build/go-journey


PKG_LIST := $(shell go list ${PKG}/... | grep -v /vendor/)

VERSION := $(shell git describe --always --long --dirty || date)
GO_FILES := $(shell find . -name '*.go' | grep -v /vendor/)

all: build

test: 
	@go test -short ${PKG_LIST}

vet:
	@go vet ${PKG_LIST}

lint:
	@for file in ${GO_FILES} ;  do \
		golint $$file ; \
	done

draft: build
	draft up

clean:
	-@rm ${OUT_BIN}

.PHONY: run build vet lint build