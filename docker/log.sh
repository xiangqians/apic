#!/bin/bash

# 查看容器日志
# -n 是 --tail 的简写形式，都是实时跟踪容器日志的最后 100 行
#docker logs -f --tail 100 apic
#docker logs -f -n 100 apic
docker logs -f apic
