MODULE_NAME:=$(shell sh -c 'cat go.mod | grep module | sed -e "s/module //"')
PROTOBUF_VERSION=$(shell sh -c ' \
	cat go.mod | grep github.com/golang/protobuf | \
	sed -e "s/[[:space:]]\+/ /g" | \
	sed -e "s/ //g" | \
	sed -e "s/github.com\/golang\/protobuf//" | \
	sed -e "s/\/\/indirect//" \
	')
GOGO_PROTOBUF_VERSION=$(shell sh -c ' \
	cat go.mod | grep github.com/gogo/protobuf | \
	sed -e "s/[[:space:]]\+/ /g" | \
	sed -e "s/ //g" | \
	sed -e "s/github.com\/gogo\/protobuf//" | \
	sed -e "s/\/\/indirect//" \
	')

.PHONY: build
all: deps install_protoc compile_proto
deps:
	go mod download
install_protoc:
	go install ${GOPATH}/pkg/mod/github.com/golang/protobuf@${PROTOBUF_VERSION}/protoc-gen-go
	go install ${GOPATH}/pkg/mod/github.com/gogo/protobuf@${GOGO_PROTOBUF_VERSION}/protoc-gen-gogofast

GOGO_ANY=Mgoogle/protobuf/any.proto=github.com/gogo/protobuf/types
GOGO_DURATION=Mgoogle/protobuf/duration.proto=github.com/gogo/protobuf/types
GOGO_STRUCT=Mgoogle/protobuf/struct.proto=github.com/gogo/protobuf/types
GOGO_TIMESTAMP=Mgoogle/protobuf/timestamp.proto=github.com/gogo/protobuf/types
GOGO_WRAPPERS=Mgoogle/protobuf/wrappers.proto=github.com/gogo/protobuf/types
GOGO_PROTOBUF_DIR=${GOPATH}/pkg/mod/github.com/gogo/protobuf@${GOGO_PROTOBUF_VERSION}
PROTOBUF_FILES_DIR=internal/pkg/broker/proto
compile_proto:
	protoc -I ${PROTOBUF_FILES_DIR} \
        -I ${GOGO_PROTOBUF_DIR} \
        --gogofast_out=${GOGO_ANY},${GOGO_DURATION},${GOGO_STRUCT},${GOGO_TIMESTAMP},${GOGO_WRAPPERS},plugins=grpc:${PROTOBUF_FILES_DIR} \
        ${PROTOBUF_FILES_DIR}/*.proto
clear:
	rm -f coverage.out coverage.html
tests: clear
	go test -covermode=count -coverprofile=coverage.out `go list ./...` | grep -q ""
	go tool cover -html=coverage.out -o coverage.html
coverage: tests
	go tool cover -func=coverage.out
format:
	go fmt `go list ./... | grep -v /vendor/`
	goimports -w -local ${MODULE_NAME} `go list -f {{.Dir}} ./...`

DOCKERFILE=build/Dockerfile
DOCKER_IMAGE=typerium
build:
	docker build -t ${DOCKER_IMAGE} -f ${DOCKERFILE} .

DEV_CMD=docker-compose -f deployments/dev.docker-compose.yml
dev_start:
	${DEV_CMD} up -d
dev_stop:
	${DEV_CMD} down