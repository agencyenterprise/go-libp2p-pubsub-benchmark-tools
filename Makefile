-include .env

all: help

LINTER := $(shell which golangci-lint 2> /dev/null)
VERSION := $(shell git describe --tags)
BUILD := $(shell git rev-parse --short HEAD)
PROJECTNAME := $(shell basename "$(PWD)")

# Go related variables.
GOBASE := $(shell pwd)
GOPATH := $(GOBASE)/vendor:$(GOBASE)
GOBIN := $(GOBASE)/.bin
GOFILES := $(wildcard *.go)

# Use linker flags to provide version/build settings
LDFLAGS=-ldflags "-X=main.Version=$(VERSION) -X=main.Build=$(BUILD)"

# Redirect error output to a file, so we can show it in development mode.
STDERR := /tmp/.$(PROJECTNAME)-stderr.txt

# PID file will keep the process id of the server
PID := /tmp/.$(PROJECTNAME).pid

# Make is verbose in Linux. Make it silent.
MAKEFLAGS += --silent

.PHONY: exists/linter
exists/linter:
ifndef LINTER
	$(error "No golangci-lint in PATH, consider doing 'GO111MODULE=on go get github.com/golangci/golangci-lint/cmd/golangci-lint@latest'")
endif

## host: run the gossip pub/sub host.
.PHONY: host
host:
	@echo "  >  Starting host..."
	@go run ./cmd/host/main.go

.PHONY: fmt
fmt:
	@echo "  >  Formatting..."
	@go list ./... | grep -v vendor | xargs gofmt -e -s -w

.PHONY: lint
lint: exists/linter
	@echo "  >  Linting..."
	@golangci-lint run ./...

## test: run all tests.
.PHONY: test
test: test/unit test/integration test/e2e
	@echo "  >  Starting tests..."
	@echo "passed all tests"

.PHONY: test/unit
test/unit:
	@echo "  >  Running unit tests..."
	@go test ./... -tags=unit -v && echo "passed unit tests"

.PHONY: test/integration
test/integration:
	@echo "  >  Running integration tests..."
	@go test ./... -tags=integration -v && echo "passed integration tests"

.PHONY: test/e2e
test/e2e:
	@echo "  >  Running e2e tests..."
	@go test ./... -tags=e2e -v && echo "passed e2e tests"

## compile: compile the binary.
compile:
	@-touch $(STDERR)
	@-rm $(STDERR)
	@-$(MAKE) -s build 2> $(STDERR)
	@cat $(STDERR) | sed -e '1s/.*/\nError:\n/'  | sed 's/make\[.*/ /' | sed "/^/s/^/     /" 1>&2

## clean: Clean build files. Runs `go clean` internally.
.PHONY: clean
clean:
	@-rm $(GOBIN)/$(PROJECTNAME) 2> /dev/null
	@-$(MAKE) go/clean

build: get fmt lint test build
	@echo "  >  Building binary..."
	@GOPATH=$(GOPATH) GOBIN=$(GOBIN) go build $(LDFLAGS) -o $(GOBIN)/$(PROJECTNAME) $(GOFILES)

.PHONY: generate
generate:
	@echo "  >  Generating dependency files..."
	@GOPATH=$(GOPATH) GOBIN=$(GOBIN) go generate $(generate)

## get: Fetch all dependencies.
.PHONY: get
get:
	@echo "  >  Checking if there is any missing dependencies..."
	@GOPATH=$(GOPATH) GOBIN=$(GOBIN) go get $(get)

.PHONY: vendor
vendor:
	@GO111MODULE=on go mod vendor

.PHONY: install
install:
	@GOPATH=$(GOPATH) GOBIN=$(GOBIN) go install $(GOFILES)

.PHONY: go/clean
go/clean:
	@echo "  >  Cleaning build cache"
	@GOPATH=$(GOPATH) GOBIN=$(GOBIN) go clean

.PHONY: kill/gopls
kill/gopls:
	@killall gopls

.PHONY: help
help: Makefile
	@echo
	@echo " Choose a command to run in "$(PROJECTNAME)":"
	@echo
	@sed -n 's/^##//p' $< | column -t -s ':' |  sed -e 's/^/ /'
	@echo
