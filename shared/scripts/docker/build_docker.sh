#!/bin/sh
set -x
set -e

env GOOS=linux go build -v ./docker/fetch.cli
docker build --force-rm -t sniperkit/dfk-fetch.cli:3.7-alpine --no-cache -f ./docker/dockerfile-alpine3.7