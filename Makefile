VERSION := $(shell echo $(shell git describe --tags) | sed 's/^v//')
COMMIT  := $(shell git log -1 --format='%H')

GOPATH := $(shell go env GOPATH)
GOBIN := $(GOPATH)/bin

all: lint install

###############################################################################
# Build / Install
###############################################################################

LD_FLAGS = -X github.com/Moonyongjung/xpla-anchor/cmd.Version=$(VERSION)

BUILD_FLAGS := -ldflags '$(LD_FLAGS)'

install: go.sum
	@echo "installing anc binary..."
	@go build -mod=readonly $(BUILD_FLAGS) -o $(GOBIN)/anc anchor.go

###############################################################################
# Tests / CI
###############################################################################

test:
	@go test -mod=readonly -race ./...
