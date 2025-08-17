#!/bin/bash

# 操作系统
os="$(go env GOOS)"
echo "OS      : ${os}"

# CPU 架构
arch="$(go env GOARCH)"
echo "ARCH    : ${arch}"

# 当前目录
cur_dir="$(cd $(dirname $0); pwd)"
echo "CUR_DIR : ${cur_dir}"

# 输出目录
out_dir="${cur_dir}/build"
echo "OUT_DIR : ${out_dir}"

# 清空输出目录
rm -rf "${out_dir}"
# 创建输出目录
mkdir -p "${out_dir}"

# 拷贝文件
cp config.ini "${out_dir}/"
cp -r example "${out_dir}/"

# 构建
echo BUILDING ...
out_name="apic-${os}-${arch}"
out_path="${out_dir}/${out_name}"
cd "${cur_dir}" && go build -ldflags="-s -w" -o "${out_path}"

# 压缩可执行文件
# $ apt install upx
#upx -9 --brute --backup "${out_path}"

# 启动命令
out_path="${out_dir}/start.sh"
cat>"${out_path}"<<EOF
#!/bin/bash
#nohup ./${out_name} >/dev/null 2>&1 &
./${out_name}
EOF
chmod +x "${out_path}"