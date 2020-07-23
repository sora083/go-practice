export GO111MODULE=on

all: practice

restapp: *.go
	go build -o practice


NAME := practice # FIXME
VERSION := 1.0.0
REVISION := $(shell git rev-parse --short HEAD)

SRCS     := $(shell find . -type f -name '*.go')
LDFLAGS  := -ldflags="-s -w -extldflags \"-static\""
NOVENDOR := $(shell go list ./... | grep -v vendor)

DIST_DIRS := find * -type d -exec

.DEFAULT_GOAL := help

ifndef GOBIN
GOBIN := $(shell echo "$${GOPATH%%:*}/bin")
endif

LINT := $(GOBIN)/golint
GOX := $(GOBIN)/gox
ARCHIVER := $(GOBIN)/archiver
GHR := $(GOBIN)/ghr

$(LINT): ; @go get github.com/golang/lint/golint
$(GOX): ; @go get github.com/mitchellh/gox
$(ARCHIVER): ; @go get github.com/mholt/archiver/cmd/arc
$(GHR): ; @go get github.com/tcnksm/ghr

.PHONY: help
help: ## Show help see: https://postd.cc/auto-documented-makefile/
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

.PHONY: deps
deps: ## Install dependency libraries
	go mod download

.PHONY: tidy
tidy: ## remove unnecessary deps
	go mod tidy

.PHONY: update-deps
update-deps: ## Update dependency libraries
	go get -u

.PHONY: build
build: deps ## Build app for linux arch
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build $(LDFLAGS) -v -o dist/$(NAME)

.PHONY: server
server: ## Run api server for test
ifeq ($(shell command -v realize 2> /dev/null),)
	go get -u github.com/oxequa/realize
endif
	realize start

# .PHONY: deploy
# deploy: build ## build and deploy module to server

# Alias
s: server
w: watch
d: deploy


# all: build

# .PHONY: clean
# clean:
# 	rm -rf practice

# .PHONY: build
# build:
# 	go build -o practice