#!/bin/bash

set -ex

echo "Starting getting dependency"

# shellcheck disable=SC2068
go get -u $@ \
  github.com/golang/protobuf/protoc-gen-go \
  google.golang.org/protobuf/proto \
  github.com/grpc-ecosystem/grpc-gateway/protoc-gen-grpc-gateway \
  github.com/grpc-ecosystem/grpc-gateway/protoc-gen-swagger

echo "Successfully getting dependency"
