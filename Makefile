.PHONY: build doc dep dev test
.DEFAULT_GOAL: build

PKG_NAME=$(shell basename `pwd`)
OPERATOR_PKGS=$(shell go list ./... | grep -v /vendor/)

# These are the values we want to pass for Version and BuildTime
VERSION=0.0.1
BUILD_TIME=$(shell date +%FT%T%z)

# This is how we want to name the binary output
SERVER_BINARY=opserver
CLIENT_BINARY=opcli

# Setup the -ldflags option for go build here, interpolate the variable values
LDFLAGS=-ldflags "-X github.com/colebrumley/operator.Version=$(VERSION) -X github.com/colebrumley/operator.BuildTime=$(BUILD_TIME)"
EXTRAFLAGS=-x -v -a -installsuffix cgo

# Listener address for the godoc server
GODOC_LISTENER=:6060

build: clean dep test
	go build -o ./bin/$(SERVER_BINARY) $(LDFLAGS) $(EXTRAFLAGS) ./server/*.go
	go build -o ./bin/$(CLIENT_BINARY) $(LDFLAGS) $(EXTRAFLAGS) ./client/*.go

doc:
	@echo $(PKG_NAME) docs available at http://localhost:6060/pkg/github.com/colebrumley/$(PKG_NAME)/
	godoc -http=$(GODOC_LISTENER)

dep:
	glide install

dev: dep
	mkdir bin 2>/dev/null; docker-compose up -d

test:
	go test $(OPERATOR_PKGS)

clean:
	@rm -Rf bin/* 2>/dev/null

clean-all:
	@rm -Rf bin/ vendor/ glide.lock