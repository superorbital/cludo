#!/usr/bin/env bash

~/bin/swagger generate server -f ./swagger.yaml --main-package=server --principal=models.ModelsPrincipal
~/bin/swagger generate client -f ./swagger.yaml
# ~/bin/swagger generate cli -f ./swagger.yaml --template-dir=./templates
