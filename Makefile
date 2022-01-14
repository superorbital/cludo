GOPATH := $(shell go env GOPATH)
GOFILES := $(wildcard *.go)

ARCH := $(shell go env GOARCH)
GITCOMMIT := $(shell git rev-parse --short=12 HEAD)
GITBRANCH := $(shell git rev-parse --abbrev-ref HEAD)
VERSION := $(shell cat VERSION)

# Use linker flags to provide version/build settings
LDFLAGS=-ldflags "$(shell go run github.com/ahmetb/govvv -flags -pkg github.com/superorbital/cludo/pkg/build)"

# Make is verbose in Linux. Make it silent.
# MAKEFLAGS += --silent

PR_NUM ?= ""

all: test build docker
.PHONY: all swagger build test clean docker docker-local-arch-build nerdctl nerdctl-local-arch-build

all-nc: test build nerdctl

swagger:
	./bin/gen-swagger.sh

# Naming the file with the os/arch makes it super simple to upload to a Github release, as is.
build:
	go mod tidy
	go run github.com/mitchellh/gox $(LDFLAGS) -output "builds/{{.OS}}_{{.Arch}}_{{.Dir}}" -osarch="darwin/amd64 darwin/arm64 linux/amd64 linux/arm64 windows/amd64" ./...

docker: docker-local-arch-build

nerdctl: nerdctl-local-arch-build

# Can't currently sideload manifest-based builds, so we
# have to split this into multiple builds.
docker-local-arch-build:
# Build images for local testing
	docker buildx build --platform linux/$(ARCH) -t superorbital/cludo:$(VERSION).git-$(GITCOMMIT)-local -t superorbital/cludo:$(VERSION)-local -t superorbital/cludo:local -f ./Dockerfile --load .
	docker buildx build --platform linux/$(ARCH) -t superorbital/cludod:$(VERSION).git-$(GITCOMMIT)-local -t superorbital/cludod:$(VERSION)-local -t superorbital/cludod:local -f ./Dockerfile.cludod --load .

nerdctl-local-arch-build:
# Build images for local testing
# Nerdctl v0.15.0 does not support multiple tags (v0.15.0)
	nerdctl build --namespace k8s.io --platform linux/$(ARCH) -t superorbital/cludo:local -f ./Dockerfile .
	nerdctl build --namespace k8s.io --platform linux/$(ARCH) -t superorbital/cludod:local -f ./Dockerfile.cludod .

docker-build-push:
ifeq ($(shell git rev-parse --abbrev-ref HEAD),main)
# This appears to be the default branch (main)
	docker buildx build --platform linux/amd64,linux/arm64 -t superorbital/cludo:$(VERSION) -t superorbital/cludo:latest -f ./Dockerfile --push .
	docker buildx build --platform linux/amd64,linux/arm64 -t superorbital/cludod:$(VERSION) -t superorbital/cludod:latest -f ./Dockerfile.cludod --push .
else ifneq ($(PR_NUM),"")
# FIXME: We appear to get through here with an empty value in some cases...
# This appears to be a Pull Request
		docker buildx build --platform linux/amd64,linux/arm64 -t superorbital/cludo:development.git-PR-$(PR_NUM) -f ./Dockerfile --push .
		docker buildx build --platform linux/amd64,linux/arm64 -t superorbital/cludod:development.git-PR-$(PR_NUM) -f ./Dockerfile.cludod --push .
else
# This appears to be a non-default branch
# (for now do nothing , as this is unusual and we don't have a cleanup strategy)
# We will also need to make GITBRANCH valid for an image tag.
		echo "Not building for non-main branch"
#docker buildx build --platform linux/amd64,linux/arm64 -t superorbital/cludo:development.git-$(GITBRANCH) -f ./Dockerfile --push .
#docker buildx build --platform linux/amd64,linux/arm64 -t superorbital/cludod:development.git-$(GITBRANCH) -f ./Dockerfile.cludod --push .
endif

test:
	go test ./...

clean:
	-rm -rf bin/*
