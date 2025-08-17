#!/bin/bash

# 运行镜像
docker run \
  -id \
  --log-opt max-size=10m \
  --log-opt max-file=1 \
  -v /opt/apic/config.ini:/app/config.ini:ro \
  -v /opt/apic/example:/app/example \
  -p 58081:58081 \
  --name apic \
  -t apic:1.0.0
# -id 交互式后台运行，-i 保持标准输入开放，-d 后台运行
# --log-opt max-size=10m 限制单个日志文件最大大小
# --log-opt max-file=1 限制日志文件数量（超出后覆盖）
# -v <宿主机路径>:<容器路径> 数据卷挂载，将宿主机目录挂载到容器路径，:ro 防止容器修改
# -p <宿主机端口>:<容器端口>：端口映射，将宿主机端口映射到容器端口
# --name test	容器命名
# -t <image_name>:<tag>：指定镜像名称和标签
