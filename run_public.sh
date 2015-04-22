#!/usr/bin/env bash

go build || exit $?
export LISTEN=127.0.0.1:1327
export ROOT=/api/auth
export PUBLIC=true

echo "LISTEN=$LISTEN"
echo "ROOT=$ROOT"
echo "PUBLIC=$PUBLIC"

./auth-server || exit $?
