#!/bin/bash
set -e

export GOPATH=$PWD:$GOPATH

current_dir=$(cd $(dirname $0)/../..; pwd)
source $current_dir/ci/tasks/docker.sh

# Setup docker environment
mount_cgroups
start_docker

cd src/github.com/apihub/apihub
make docker-test


