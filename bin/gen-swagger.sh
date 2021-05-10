#!/usr/bin/env bash

# $1 - [opt] TEMPLATE_DIR / "./templates"

set -euxo pipefail

if [[ -z ${1+x} ]]; then
  export TEMPLATE_DIR="./templates"
else
  export TEMPLATE_DIR="$4"
fi

echo "[INFO} TEMPLATE_DIR: ${TEMPLATE_DIR}"

commands=( "swagger" )
sites=("https://github.com/go-swagger/go-swagger#installing")

for i in "${!commands[@]}"; do
  if ! command -v "${commands[$i]}" &> /dev/null; then
  echo "[ERROR] ${commands[$i]} (${sites[$i]}) must be installed."
  exit 1
fi
done

DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" >/dev/null 2>&1 && pwd )"

cd "${DIR}/.."

~/bin/swagger generate server -f ./swagger.yaml --main-package=cludod --principal=models.ModelsPrincipal --template-dir="${TEMPLATE_DIR}"
~/bin/swagger generate client -f ./swagger.yaml --template-dir="${TEMPLATE_DIR}"
# ~/bin/swagger generate cli -f ./swagger.yaml --cli-app-name=cludo-api --template-dir="${TEMPLATE_DIR}"

go get -u ./...
