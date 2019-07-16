#!/usr/bin/env bash

GOGO_PROTOBUF_VERSION=$(cat go.mod | grep github.com/gogo/protobuf | sed -e "s/[[:space:]]\+/ /g" \
    | sed -e "s/ //g" | sed -e "s/github.com\/gogo\/protobuf//" | sed -e "s/\/\/indirect//")
GOGO_PROTOBUF_DIR=$GOPATH/pkg/mod/github.com/gogo/protobuf@${GOGO_PROTOBUF_VERSION}

GOGO_ANY=Mgoogle/protobuf/any.proto=github.com/gogo/protobuf/types
GOGO_DURATION=Mgoogle/protobuf/duration.proto=github.com/gogo/protobuf/types
GOGO_STRUCT=Mgoogle/protobuf/struct.proto=github.com/gogo/protobuf/types
GOGO_TIMESTAMP=Mgoogle/protobuf/timestamp.proto=github.com/gogo/protobuf/types
GOGO_WRAPPERS=Mgoogle/protobuf/wrappers.proto=github.com/gogo/protobuf/types

PROTOBUF_FILES_DIR=internal/pkg/broker/proto

protoc -I .. \
    -I ${PROTOBUF_FILES_DIR} \
    -I ${GOGO_PROTOBUF_DIR} \
    --gogofast_out=${GOGO_ANY},${GOGO_DURATION},${GOGO_STRUCT},${GOGO_TIMESTAMP},${GOGO_WRAPPERS},plugins=grpc:${PROTOBUF_FILES_DIR} \
    ${PROTOBUF_FILES_DIR}/*.proto