#!/usr/bin/env sh

set -e
set -o pipefail

if [ -n "${CLUDOD_CONFIG}" ]; then
    mkdir -p /etc/cludod
    echo "${CLUDOD_CONFIG}" > /etc/cludod/cludod.yaml
fi

if [ -z "${PORT}" ]; then
    export PORT=80
fi

if [ -z "${@}" ]; then
    cludod --port=${PORT} --scheme=http --host=0.0.0.0
else
    ${@}
fi
