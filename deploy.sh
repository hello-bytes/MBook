#!/bin/bash

if [ ! -f deploy.sh ]; then
    echo 'deploy.sh must be run within its container folder' 1>&2
    exit 1
fi

function build(){
    DIRPWD=`pwd`
    export GOPATH=$DIRPWD

    rm -rf build/v1.0.3
    mkdir build
    mkdir build/v1.0.3

    # build go binary
    cd src/main
    GOOS=linux go build main.go website.go

    # copy and delete go binary
    mv main ../../build/v1.0.3/xbook
    rm main

    # copy resource files
    cd ../..
    mkdir build/v1.0.3/resource
    cp -R resource/* build/v1.0.3/resource/

    rm build/v1.0.3/resource/config/env.config

    cd build/v1.0.3
    zip -q -r xbook_v1.0.3.zip ./*
}

function buildfortest(){
    DIRPWD=`pwd`
    export GOPATH=$DIRPWD

    # build go binary
    cd src/main
    go build main.go website.go

    # copy and delete go binary
    cp main ../../mbook
    rm main

    cd ../..
}

function buildmac(){
    DIRPWD=`pwd`
    export GOPATH=$DIRPWD

    rm -rf build/v1.0.3
    mkdir build
    mkdir build/v1.0.3

    # build go binary
    cd src/main
    go build main.go website.go

    # copy and delete go binary
    mv main ../../build/v1.0.3/xbook
    rm main

    # copy resource files
    cd ../..
    mkdir build/v1.0.3/resource
    cp -R resource/* build/v1.0.3/resource/

    rm build/v1.0.3/resource/config/env.config

    cd build/v1.0.3
    zip -q -r xbook_v1.0.3.zip ./*
}

function start(){
    echo "starting"
    basepath=$(cd `dirname $0`; pwd)
    echo ${basepath}
    eval "cd ${basepath} && nohup ./xbook console 2>&1 &"
}

function stop(){
    echo "stoping"
    PROCESS=`ps -ef|grep ./xbook|grep -v grep|grep -v PPID|awk '{ print $2}'`
    for i in $PROCESS
    do
        echo "Kill the xbook process [ $i ]"
        kill -9 $i
    done
    # eval "ps -ef | grep procedure_name | grep -v grep | awk '{print "./xbook"}' | xargs kill -9"
}

function restart(){
    stop
    start
}

function deploy(){
    basepath=$(cd `dirname $0`; pwd)

    cd ${basepath}
    filename=$1
    echo $filename

    # unzip ${filename} 
    echo "start unzip file ${filename}"
    unzip -q -d ./xbook_temp ${filename}
    echo "unzip finish."

    # copy executable
    echo "delete and copy xbook to root path"
    rm xbook
    cp xbook_temp/xbook ./
    chmod +x xbook

    echo "copy pubilc and templates..."
    mkdir resource
    rm -rf resource/public
    mkdir resource/public
    cp -R xbook_temp/resource/public/* resource/public/

    rm -rf resource/templates
    mkdir resource/templates
    cp -R xbook_temp/resource/templates/* resource/templates/

    echo "remove the temp(unzip) file..."
    rm -rf xbook_temp/
    echo "done!!!"
}

function help(){
    echo "Usage : "
    echo "     bash deploy.sh build"
    echo "     e.g : bash deploy.sh build"
}

case $1 in
    deploy)
        deploy $2
        ;;
    build)
        build
        ;;
    buildmac)
        buildmac
        ;;
    debug)
        buildfortest
        ;;
    start)
        start
        ;;
    stop)
        stop
        ;;
    restart)
        restart
        ;;
    help)
        help
        exit 1;;
    ?)
        help
        exit 1;;
esac
