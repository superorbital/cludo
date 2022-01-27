#!/usr/bin/env bash

set -euo pipefail

cp ~/.cludo/cludo.yaml ~/.cludo/cludo.yaml.orig

# trap ctrl-c and call ctrl_c()
trap cleanup INT EXIT

function cleanup() {
  cp ~/.cludo/cludo.yaml.orig ~/.cludo/cludo.yaml
}

# Fix this once custpm config paths work again
# https://github.com/superorbital/cludo/issues/108
for i in $(ls -C1 ~/.cludo/cludo-test-*); do
  echo -e "\n\n"
  cp $i ~/.cludo/cludo.yaml
  grep ssh_key_paths ~/.cludo/cludo.yaml | awk -F"[" '{print $2}' | awk -F"]" '{print $1}'
  ./builds/darwin_amd64_cludo exec aws sts get-caller-identity
  cp ~/.cludo/cludo.yaml.orig ~/.cludo/cludo.yaml
done

exit 0

