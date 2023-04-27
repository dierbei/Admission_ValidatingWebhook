FROM golang:1.18-alpine AS builder

# ENV 设置环境变量
ENV GO111MODULE=on
ENV GOPROXY=https://goproxy.io,direct

# COPY 源路径 目标路径
COPY . /app

# 设置工作目录
WORKDIR /app

# RUN 执行 go build .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o main .

# FROM 基于 alpine:latest
FROM alpine:latest

# RUN 设置代理镜像
RUN echo -e http://mirrors.ustc.edu.cn/alpine/v3.13/main/ > /etc/apk/repositories

# RUN 设置 Asia/Shanghai 时区
RUN apk --no-cache add tzdata  && \
    ln -sf /usr/share/zoneinfo/Asia/Shanghai /etc/localtime && \
    echo "Asia/Shanghai" > /etc/timezone

# COPY 源路径 目标路径 从镜像中 COPY
COPY --from=builder /app/main /opt/main
COPY --from=builder /app/cert /opt/cert

# WORKDIR 设置工作目录
WORKDIR /opt

# CMD 设置启动命令
CMD ["./main"]
