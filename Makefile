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
	./bin/gen-swagger.sh

# Naming the file with the os/arch makes it super simple to upload to a Github release, as is.
build:
	go get github.com/ahmetb/govvv
	go mod tidy
	gox $(LDFLAGS) -output "builds/{{.OS}}_{{.Arch}}_{{.Dir}}" -osarch="darwin/amd64 darwin/arm64 linux/amd64 linux/arm64 windows/amd64" ./...

docker: docker-build docker-tag

# Can't currently sideload manifest-based builds, so we
# have to split this into multiple builds.
docker-build:
	docker buildx build --platform linux/amd64 -t superorbital/cludo:$(VERSION).git-$(GITCOMMIT) -t superorbital/cludo:$(VERSION) -t superorbital/cludo:latest -f ./Dockerfile --load .
	docker buildx build --platform linux/arm64 -t superorbital/cludo:$(VERSION).git-$(GITCOMMIT) -t superorbital/cludo:$(VERSION) -t superorbital/cludo:arm64-latest -f ./Dockerfile --load .
	docker buildx build --platform linux/amd64 -t superorbital/cludod:$(VERSION).git-$(GITCOMMIT) -t superorbital/cludod:$(VERSION) -t superorbital/cludod:latest -f ./Dockerfile.cludod --load .
	docker buildx build --platform linux/arm64 -t superorbital/cludod:$(VERSION).git-$(GITCOMMIT) -t superorbital/cludod:$(VERSION) -t superorbital/cludod:arm64-latest -f ./Dockerfile.cludod --load .

docker-push:
	docker buildx build --platform linux/amd64,linux/arm64 -t superorbital/cludo:$(VERSION).git-$(GITCOMMIT) -t superorbital/cludo:$(VERSION) -t superorbital/cludo:latest -f ./Dockerfile --push .
	docker buildx build --platform linux/amd64,linux/arm64 -t superorbital/cludod:$(VERSION).git-$(GITCOMMIT) -t superorbital/cludod:$(VERSION) -t superorbital/cludod:latest -f ./Dockerfile.cludod --push .

test:
	go test ./...

clean:
	-rm -rf bin/*
