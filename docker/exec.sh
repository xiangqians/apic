#!/bin/bash

# 以 root 用户进入容器

# 容器内部已安装 /bin/bash
#docker exec -it -u root apic /bin/bash

#docker exec -it -u root apic sh
docker exec -it -u root apic env ENV=/root/.profile sh -c "echo \"alias ll='ls -alF'\" > /root/.profile && sh"
# ENV 变量指定了 sh 启动时要加载的配置文件
# 设置 ll 命令
# $ echo "alias ll='ls -alF'" >> /etc/profile && source /etc/profile
# -a：显示隐藏文件（如 .bashrc）
# -l：长格式显示（权限、所有者、大小等）
# -F：标记文件类型（如 / 表示目录，* 表示可执行文件）
