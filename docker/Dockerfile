# [构建]

# 基于 Alpine Linux 的官方 Go 1.24 镜像，包含完整的 Go 编译环境
FROM golang:1.24-alpine AS builder

# 设置工作目录
WORKDIR /app

# 复制源码
COPY ./ ./

# 设置代理
RUN go env -w GOPROXY=https://goproxy.cn,direct

# 构建
RUN go build -ldflags="-s -w" -o apic

# 压缩可执行文件
#RUN apk update
#RUN apk add --no-cache upx
#RUN upx -9 --brute --backup apic


# [镜像]

# 基于 Alpine Linux 发行版的最小化镜像，不包含非必要组件
FROM alpine:3.20

# 设置工作目录
WORKDIR /app

# 从 builder 阶段文件
COPY --from=builder /app/apic ./
COPY --from=builder /app/config.ini ./
COPY --from=builder /app/example ./example

# 设置启动命令
CMD [ "/app/apic" ]