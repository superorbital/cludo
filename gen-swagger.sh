#!/usr/bin/env bash

go-swagger generate server -f ./swagger.yaml --main-package=server
go-swagger generate client -f ./swagger.yaml
go-swagger generate cli -f ./swagger.yaml
