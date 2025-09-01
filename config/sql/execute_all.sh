#!/bin/bash

# ProRent 数据库初始化脚本
# 执行所有 SQL 文件，按顺序初始化数据库

echo "=== ProRent 数据库初始化脚本 ==="
echo "开始时间: $(date)"
echo ""

# 数据库连接信息 (从 settings.yml 读取)
DB_HOST="127.0.0.1"
DB_PORT="3306"
DB_USER="root"
DB_PASS="123456"
DB_NAME="rentpro_admin"

# 检查 MySQL 客户端
if ! command -v mysql &> /dev/null; then
    echo "❌ 错误: MySQL 客户端未安装"
    exit 1
fi

# 测试数据库连接
echo "🔍 测试数据库连接..."
if ! mysql -h"$DB_HOST" -P"$DB_PORT" -u"$DB_USER" -p"$DB_PASS" -e "USE $DB_NAME;" 2>/dev/null; then
    echo "❌ 错误: 无法连接到数据库 $DB_NAME"
    echo "请检查数据库连接信息:"
    echo "  主机: $DB_HOST:$DB_PORT"
    echo "  用户: $DB_USER"
    echo "  数据库: $DB_NAME"
    exit 1
fi
echo "✅ 数据库连接成功"
echo ""

# 执行初始化文件
echo "📋 执行初始化文件..."
mysql -h"$DB_HOST" -P"$DB_PORT" -u"$DB_USER" -p"$DB_PASS" "$DB_NAME" < init/01-mysql-begin.sql
if [ $? -eq 0 ]; then
    echo "✅ 执行: 01-mysql-begin.sql"
else
    echo "❌ 执行失败: 01-mysql-begin.sql"
    exit 1
fi

# 执行数据文件
echo ""
echo "📊 执行数据文件..."
for file in data/*.sql; do
    if [ -f "$file" ]; then
        echo "执行: $(basename "$file")"
        mysql -h"$DB_HOST" -P"$DB_PORT" -u"$DB_USER" -p"$DB_PASS" "$DB_NAME" < "$file"
        if [ $? -eq 0 ]; then
            echo "✅ 成功: $(basename "$file")"
        else
            echo "❌ 失败: $(basename "$file")"
            exit 1
        fi
    fi
done

# 执行结束文件
echo ""
echo "📋 执行结束文件..."
mysql -h"$DB_HOST" -P"$DB_PORT" -u"$DB_USER" -p"$DB_PASS" "$DB_NAME" < init/02-mysql-end.sql
if [ $? -eq 0 ]; then
    echo "✅ 执行: 02-mysql-end.sql"
else
    echo "❌ 执行失败: 02-mysql-end.sql"
    exit 1
fi

echo ""
echo "🎉 数据库初始化完成!"
echo "结束时间: $(date)"
echo ""
echo "📝 下一步操作:"
echo "1. 检查数据库表是否创建成功"
echo "2. 验证基础数据是否插入成功"
echo "3. 运行 Go 迁移文件 (如果需要)"
echo "4. 启动应用程序"
