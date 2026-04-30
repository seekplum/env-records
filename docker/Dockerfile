ARG GO_VERSION=${GO_VERSION:-"1.25-alpine3.22"}

# 当前使用go run app.go直接运行。docker相关配置暂时无用。
FROM golang:${GO_VERSION} AS builder

ADD . /go/src/app
WORKDIR /go/src/app
RUN BUILD_OUTPUT=/env-records ./bin/build.sh local_build

FROM alpine:latest

WORKDIR /
RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.ustc.edu.cn/g' /etc/apk/repositories && apk add --no-cache tzdata
ENV TZ=Asia/Shanghai
COPY --from=builder /env-records /
EXPOSE 10000
ENTRYPOINT [ "/env-records" ]
