#!/usr/bin/env bash

go build || exit $?
export LISTEN=127.0.0.1:1324
export ROOT=/api/auth

echo "LISTEN=$LISTEN"
echo "ROOT=$ROOT"
echo "PUBLIC=false"

./auth-server || exit $?
