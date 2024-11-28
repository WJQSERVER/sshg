#!/bin/bash
# Copyright (c) 2024, WJQSERVER STUDIO. Follow the WSL License.

# 获取配置信息
read -p "请输入SSHG的端口号(默认22): " port

# 创建目录
mkdir /etc/sshg
touch /etc/sshg/config.toml
