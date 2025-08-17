#!/bin/bash

# 停止容器
docker stop apic

# 删除容器
docker rm apic

# 删除镜像
docker rmi apic:1.0.0
