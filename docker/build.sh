#!/bin/bash

cd "$(dirname $0)/.."

# 构建镜像
docker build \
  -f docker/Dockerfile \
  -t apic:1.0.0 .
# -f /path/to/Dockerfile：指定 Dockerfile 文件路径
# -t <image_name>:<tag>：指定镜像名称和标签
# .：表示当前目录作为构建的上下文
