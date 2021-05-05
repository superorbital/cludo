#!/usr/bin/env bash

./swagger generate server -f ./swagger.yaml --main-package=server --flag-strategy=pflag --exclude-main
./swagger generate client -f ./swagger.yaml
# ./swagger generate cli -f ./swagger.yaml --template-dir=./templates
