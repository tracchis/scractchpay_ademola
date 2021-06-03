.PHONY: all clean build test test-unit test-integration

TIMESTAMP=$(shell date +"%Y%m%d.%H%M")
BUILDSTAMP := $(TIMESTAMP)

GIT_HASHID ?= $(shell git rev-parse --short HEAD 2> /dev/null )
ifneq (,$(GIT_HASHID))
  BUILDSTAMP := $(BUILDSTAMP)~$(GIT_HASHID)
endif

VERSION ?= dev
ifneq (,$(findstring PR-,$(VERSION))$(findstring dev,$(VERSION)))
  BUILDSTAMP := $(BUILDSTAMP)~$(VERSION)

  GIT_TAG ?= $(shell git describe --tag --abbrev=0 2> /dev/null )
  ifeq (,$(GIT_TAG))
    GIT_TAG=0.0.0
  endif

  VERSION=$(GIT_TAG)
endif

GOOS ?= linux
GOARCH ?= amd64
BUILD_DIR ?= bin/$(GOOS).$(GOARCH)
BINARY=$(BUILD_DIR)/$(NAME)
BUILD_FLAGS=-ldflags="-s -w -X main.Version=$(VERSION) -X main.Buildstamp=$(BUILDSTAMP)"

all: build

clean:
	rm -rf bin/

build:
	go build -v $(BUILD_FLAGS) -o "$(BINARY)" cmd/main.go cmd/routes.go cmd/config.go

test: test-unit

test-unit:
	go test -v -race -tags=unit -coverprofile=coverage.txt -covermode=atomic ./...

test-integration:
	@(grep -rl "^// +build integration" . > /dev/null) || (echo "no integration tests" 1>&2; exit 1)
	go test -v -race -tags=integration ./...