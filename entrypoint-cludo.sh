#!/usr/bin/env sh

set -e
set -o pipefail

if [ -n "${CLUDO_CONFIG}" ]; then
    mkdir -p /etc/cludo
    echo "${CLUDO_CONFIG}" > /etc/cludo/cludo.yaml
fi

if [ -z "${@}" ]; then
    cludo --help
else
    ${@}
fi
