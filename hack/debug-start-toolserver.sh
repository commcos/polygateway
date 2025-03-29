#!/usr/bin/env bash

unameOut="$(uname -s)"
case "${unameOut}" in
    Linux*)     machine=Linux;;
    Darwin*)    machine=Mac;;
    CYGWIN*)    machine=Cygwin;;
    MINGW*)     machine=MinGw;;
    *)          machine="UNKNOWN:${unameOut}"
esac

READLINK=readlink
if [ "$machine" = "Mac" ];then
  READLINK=greadlink
fi

CURRENT_DIR=`dirname $($READLINK -f $0)`
PROJECT_ROOT=${CURRENT_DIR%%/hack}


go run cmd/tools/tools.go --v=5

