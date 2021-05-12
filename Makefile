GOPATH := $(shell go env GOPATH)
GOFILES := $(wildcard *.go)

# Use linker flags to provide version/build settings
LDFLAGS=-ldflags "$(shell govvv -flags -pkg github.com/superorbital/cludo/pkg/build)"

# Make is verbose in Linux. Make it silent.
# MAKEFLAGS += --silent

all: test build
.PHONY: all swagger build test clean

swagger:
	./gen-swagger.sh

build:
	gox $(LDFLAGS) -output "bin/{{.OS}}/{{.Arch}}/{{.Dir}}" -osarch !darwin/386 ./...

test:
	go test ./...

clean:
	-rm -rf bin/*
