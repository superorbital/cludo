#!/usr/bin/env bash

docker run --rm -it \
    --user $(id -u):$(id -g) \
    -e GOPATH=$HOME/go:/go \
    -v $HOME:$HOME \
    -w $(pwd) \
    quay.io/goswagger/swagger \
    generate server -f ./swagger.yaml --main-package=server --template-dir=/Users/robert/src/go-swagger/generator/templates
    # --exclude-main

docker run --rm -it \
    --user $(id -u):$(id -g) \
    -e GOPATH=$HOME/go:/go \
    -v $HOME:$HOME \
    -w $(pwd) \
    quay.io/goswagger/swagger \
    generate client -f ./swagger.yaml --template-dir=/Users/robert/src/go-swagger/generator/templates

docker run --rm -it \
    --user $(id -u):$(id -g) \
    -e GOPATH=$HOME/go:/go \
    -v $HOME:$HOME \
    -w $(pwd) \
    quay.io/goswagger/swagger \
    generate cli -f ./swagger.yaml --template-dir=/Users/robert/src/go-swagger/generator/templates
