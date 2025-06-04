#!/bin/bash

# 获取命令行参数
DOMAIN=$1

# 检查是否提供域名参数
if [ -z "$DOMAIN" ]; then
  echo "用法: $0 example.com"
  exit 1
fi

# 定义要查询的记录类型
RECORD_TYPES=("A" "AAAA" "CNAME" "MX" "NS" "TXT" "SOA")

echo "🔍 正在查询域名: $DOMAIN"
echo "=============================="

# 循环查询每种记录
for TYPE in "${RECORD_TYPES[@]}"; do
  echo -e "\n📌 $TYPE 记录:"
  dig +short "$DOMAIN" "$TYPE"
done

echo -e "\n✅ 查询完成。"

