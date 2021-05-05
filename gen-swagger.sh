#!/usr/bin/env bash

./swagger generate server -f ./swagger.yaml --main-package=server --template-dir=./templates
./swagger generate client -f ./swagger.yaml --template-dir=./templates
./swagger generate cli -f ./swagger.yaml --template-dir=./templates
