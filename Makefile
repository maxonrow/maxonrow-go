DEP := $(shell command -v dep 2> /dev/null)
PACKAGES=$(shell go list ./... | grep -v '/vendor/')
LDFLAGS= -ldflags "-X github.com/maxonrow/maxonrow-go/version.GitCommit=`git rev-parse --short=8 HEAD`"
TAGS=-tags 'maxonrow'

export GO111MODULE = on

all: deps build install test

deps:
	go mod vendor

build:
	go build $(LDFLAGS) $(TAGS) -mod vendor -o ./build/mxwd ./cmd/mxwd
	#go build $(LDFLAGS) $(TAGS) -mod vendor -o ./build/mxwcli ./cmd/mxwcli


install:
	go install $(LDFLAGS) $(TAGS) -mod vendor ./cmd/mxwd
	#go install $(LDFLAGS) $(TAGS) -mod vendor ./cmd/mxwcli

test:
	go test $(PACKAGES)


.PHONY: all desp build install test