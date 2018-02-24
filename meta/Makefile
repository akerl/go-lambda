.PHONY: default local build release manual clean lint vet fmt test deps init update

HELPER_PATH := $(dir $(realpath $(lastword $(MAKEFILE_LIST))))

include Makefile.local
VERSION ?= $(shell git describe --tags --always --dirty --match=v* 2>/dev/null)
export GOPATH = $(CURDIR)/.gopath
BIN = $(GOPATH)/bin
BASEDIR = $(GOPATH)/src/$(NAMESPACE)
BASE = $(BASEDIR)/$(PACKAGE)
GOFILES = $(shell find . -type f -name '*.go' ! -path './.*' ! -path './vendor/*')
OSLIST ?= linux darwin

GO = go
GOFMT = gofmt
GOX = $(BIN)/gox
GOLINT = $(BIN)/golint
GODEP = $(BIN)/dep

default: build

local: $(BASE) deps $(GOX) fmt lint vet test
	cd $(BASE) && $(GOX) \
		-ldflags '-X $(NAMESPACE)/$(PACKAGE)/cmd.Version=$(VERSION)' \
		-gocmd="$(GO)" \
		-output="bin/$(PACKAGE)_{{.OS}}" \
		-os="$(OSLIST)" \
		-arch="amd64"
	@echo "Build completed"

clean:
	rm -rf $(GOPATH) bin

lint: $(GOLINT)
	$(GOLINT) -set_exit_status $$($(GO) list -f '{{.Dir}}' ./...)

vet:
	cd $(BASE) && $(GO) vet ./...

fmt:
	@echo "Running gofmt on $(GOFILES)"
	@files=$$($(GOFMT) -l $(GOFILES)); if [ -n "$$files" ]; then \
		  echo "Error: '$(GOFMT)' needs to be run on:"; \
		  echo "$${files}"; \
		  exit 1; \
		  fi;

test: deps
	cd $(BASE) && $(GO) test ./...

init: $(BASE) $(GODEP)
	cd $(BASE) && $(GODEP) init

update: $(BASE) $(GODEP)
	cd $(BASE) && $(GODEP) ensure -update

deps: $(BASE) $(GODEP)
	cd $(BASE) && $(GODEP) ensure

$(BASEDIR):
	mkdir -p $(BASEDIR)

$(BASE): $(BASEDIR)
	ln -s $(CURDIR) $(BASE)

$(GOLINT): $(BASE)
	$(GO) get github.com/golang/lint/golint

$(GOX): $(BASE)
	$(GO) get github.com/mitchellh/gox

$(GODEP): $(BASE)
	$(GO) get github.com/golang/dep/cmd/dep

PKGFORGE_MAKE = make -f $(HELPER_PATH)/pkgforge-helper/Makefile

manual:
	$(PKGFORGE_MAKE) manual

build:
	$(PKGFORGE_MAKE)

release:
	$(PKGFORGE_MAKE) release

