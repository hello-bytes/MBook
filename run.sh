#!/usr/bin/env bash

set -e

if [ ! -f run.sh ]; then
    echo 'run.sh must be run within its container folder' 1>&2
    exit 1
fi

if [ ! -d runtime ]; then
    mkdir runtime
fi

if [ ! -d runtime/log ]; then
	mkdir runtime/log
fi

if [ ! -d runtime/pid ]; then
	mkdir runtime/pid
fi

DIRPWD=`pwd`
export GOPATH=$DIRPWD

cd src/main
go run main.go website.go console
#go build