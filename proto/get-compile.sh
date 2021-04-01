#!/bin/bash

echo "Starting compile proto"

SRC=$(realpath $(cd -P "$( dirname "${BASH_SOURCE[0]}" )" && pwd ))

GITVER=$(git describe --exact-match 2> /dev/null || echo "`git symbolic-ref HEAD 2> /dev/null | cut -b 12-`-`git log --pretty=format:\"%h\" -1`")

set -ve

./get-googleapi.sh

pushd $SRC &> /dev/null

#export PATH=$PATH:$SRC/node_modules/.bin
export PATH=$PATH:$GOPATH/bin

for VER in v*; do
  pushd $VER &> /dev/null

  for i in $(find ./ -mindepth 1 -maxdepth 1 -type d); do
    SERVICE=$(basename $i)

    pushd $SRC/$VER/$SERVICE &> /dev/null

    protoc \
      -I. \
      -I/usr/local/include \
      -I${GOPATH}/src \
      -I${GOPATH}/src/belajar/reservation/proto \
      -I${GOPATH}/src/belajar/reservation/proto/lib \
      --go_out=plugins=grpc:$GOPATH/src \
      --swagger_out=logtostderr=true:$GOPATH/src/belajar/reservation/swagger \
      --grpc-gateway_out=logtostderr=true:$GOPATH/src \
      *.proto

    popd &> /dev/null
  done

  popd &> /dev/null
done

# remove dopinang-och-go.json and generate merge swagger json file and remove each raw files
rm -rf $GOPATH/src/belajar/reservation/swagger/docs.json
jq -s 'reduce .[] as $item ({}; . * $item)|del(.info)|del(..|.tags?)' \
$GOPATH/src/belajar/reservation/swagger/*swagger.json >> $GOPATH/src/belajar/reservation/swagger/docs.json
rm -rf $GOPATH/src/belajar/reservation/swagger/*swagger.json


# build/install go output
for VER in v*; do
  pushd $VER &> /dev/null

  for i in $(find ./ -mindepth 1 -maxdepth 1 -type d); do
    SERVICE=$(basename $i)

    pushd $SRC/$VER/$SERVICE &> /dev/null
    go install
    popd &> /dev/null
  done

  popd &> /dev/null
done

popd &> /dev/null

echo "Successfully compile proto"