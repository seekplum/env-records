#!/bin/sh

ROOT="$(cd "$(dirname "$(dirname "$0")")" && pwd)"
cd "${ROOT}"

source ./bin/constants.sh

BUILD_OUTPUT=${BUILD_OUTPUT:-"build/env-records"}

function docker_build() {
    docker run --rm \
        -v ${ROOT}:/code \
        -w /code \
        golang:${GO_VERSION} /bin/sh /code/bin/build.sh local_build
}

function local_build() {
    echo "Building to ${BUILD_OUTPUT}"
    # go 文件需要配置标签
    # //go:build test
    # // +build test
    go env -w GOPROXY=https://goproxy.cn,direct && CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -tags "!test" -o ${BUILD_OUTPUT} .
}

if [ "$1" = "local_build" ]; then
	local_build
else
	docker_build
fi
