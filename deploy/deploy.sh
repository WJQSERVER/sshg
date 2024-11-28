#!/bin/bash
# Copyright (c) 2024, WJQSERVER STUDIO. Follow the WSL License.

# install packages
install() {
    if [ $# -eq 0 ]; then
        echo "ARGS NOT FOUND"
        return 1
    fi

    for package in "$@"; do
        if ! command -v "$package" &>/dev/null; then
            if command -v dnf &>/dev/null; then
                dnf -y update && dnf install -y "$package"
            elif command -v yum &>/dev/null; then
                yum -y update && yum -y install "$package"
            elif command -v apt &>/dev/null; then
                apt update -y && apt install -y "$package"
            elif command -v apk &>/dev/null; then
                apk update && apk add "$package"
            else
                echo "UNKNOWN PACKAGE MANAGER"
                return 1
            fi
        fi
    done

    return 0
}

# 检查是否为root用户
if [ "$EUID" -ne 0 ]; then
    echo "请以root用户运行此脚本"
    exit 1
fi

# 安装依赖包
install curl wget sed

# 查看当前架构是否为linux/amd64或linux/arm64
ARCH=$(uname -m)
if [ "$ARCH" != "x86_64" ] && [ "$ARCH" != "aarch64" ]; then
    echo " $ARCH 架构不被支持"
    exit 1
fi

# 重写架构值,改为amd64或arm64
if [ "$ARCH" == "x86_64" ]; then
    ARCH="amd64"
elif [ "$ARCH" == "aarch64" ]; then
    ARCH="arm64"
fi

# 获取配置信息
read -p "请输入SSHG的端口号(默认22): " port
if [ -z "$port" ]; then
    port=22
fi

read -p "请输入主机名: " hostname
if [ -z "$hostname" ]; then
    # 获取主机名
    hostname=$(hostname)
fi


read -p "请输入Telegram Bot Token: " token
if [ -z "$token" ]; then
    echo "Telegram Bot Token不能为空"
    exit 1
fi

read -p "请输入Telegram Chat ID: " chat_id
if [ -z "$chat_id" ]; then
    echo "Telegram Chat ID不能为空"
    exit 1
fi

# 创建目录
mkdir /etc/sshg
touch /etc/sshg/config.toml
mkdir /usr/local/sshg

# 构造config.toml文件
cat << EOF > /etc/sshg/config.toml
[server]
hostname = "$hostname"

[tgbot]
token = "$token"
chatid = $chat_id

[log]
logfilepath = "/var/log/sshg.log" 
maxlogsize = 5 # MB

EOF

# 构造systemd服务文件
cat << EOF > /etc/systemd/system/sshg.service
[Unit]
Description=SSH Guard Service
After=network.target

[Service]
ExecStart=/bin/bash -c '/usr/local/sshg/sshg -cfg /etc/sshg/config.toml'
WorkingDirectory=/usr/local/sshg
Restart=always
User=root
Group=root

[Install]
WantedBy=multi-user.target

EOF


# 获取最新版本号
VERSION=$(curl -s https://raw.githubusercontent.com/WJQSERVER/sshg/main/VERSION)
wget -q -O /usr/local/sshg/VERSION https://raw.githubusercontent.com/WJQSERVER/sshg/main/VERSION

# 拉取最新版本的SSHG
wget https://github.com/WJQSERVER/sshg/releases/download/${VERSION}/sshg-linux-$ARCH.tar.gz -O /usr/local/sshg.tar.gz
# 解压SSHG
install tar
# 解压SSHG
tar -zxvf /usr/local/sshg.tar.gz -C /usr/local/sshg
# 修改文件名
mv /usr/local/sshg/sshg-linux-$ARCH /usr/local/sshg/sshg
# 删除压缩包
rm /usr/local/sshg.tar.gz
# 赋予执行权限
chmod +x /usr/local/sshg/sshg

# 启动SSHG
systemctl start sshg.service
# 设置开机启动
systemctl enable sshg.service

echo "SSHG安装成功"