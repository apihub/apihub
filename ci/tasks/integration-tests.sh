#!/bin/bash
set -e

export GOPATH=$PWD:$GOPATH

cd src/github.com/apihub/apihub
ginkgo -r -p -race integration
