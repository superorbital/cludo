GOPATH := $(shell go env GOPATH)
GOFILES := $(wildcard *.go)

GITCOMMIT := $(shell git rev-parse --short HEAD)
VERSION := $(shell cat VERSION)

# Use linker flags to provide version/build settings
LDFLAGS=-ldflags "$(shell govvv -flags -pkg github.com/superorbital/cludo/pkg/build)"

# Make is verbose in Linux. Make it silent.
# MAKEFLAGS += --silent

all: test build docker
.PHONY: all swagger build test clean docker docker-build docker-tag docker-push

swagger:
	./gen-swagger.sh

build:
	gox $(LDFLAGS) -output "bin/{{.OS}}/{{.Arch}}/{{.Dir}}" -osarch !darwin/386 ./...

docker: docker-build docker-tag

docker-build:
	docker build -f Dockerfile -t superorbital/cludo:$(VERSION).git-$(GITCOMMIT) .
	docker build -f Dockerfile.cludod -t superorbital/cludod:$(VERSION).git-$(GITCOMMIT) .

docker-tag:
	docker tag superorbital/cludo:$(VERSION).git-$(GITCOMMIT) superorbital/cludo:$(VERSION)
	docker tag superorbital/cludod:$(VERSION).git-$(GITCOMMIT) superorbital/cludod:$(VERSION)
	docker tag superorbital/cludo:$(VERSION).git-$(GITCOMMIT) superorbital/cludo:latest
	docker tag superorbital/cludod:$(VERSION).git-$(GITCOMMIT) superorbital/cludod:latest

docker-push:
	docker push superorbital/cludo:$(VERSION).git-$(GITCOMMIT)
	docker push superorbital/cludo:latest
	docker push superorbital/cludod:$(VERSION).git-$(GITCOMMIT)
	docker push superorbital/cludod:latest

test:
	go test ./...

clean:
	-rm -rf bin/*
