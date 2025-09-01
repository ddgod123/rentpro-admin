#!/bin/bash

# ProRent 数据库重置脚本
# 删除所有表并重新初始化数据库

echo "=== ProRent 数据库重置脚本 ==="
echo "开始时间: $(date)"
echo "⚠️  警告: 此操作将删除数据库中的所有表和数据!"
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

# 确认操作
read -p "确定要删除数据库 $DB_NAME 中的所有表吗? (y/N): " -n 1 -r
echo
if [[ ! $REPLY =~ ^[Yy]$ ]]; then
    echo "❌ 操作已取消"
    exit 1
fi

# 删除所有表 - 使用更简单的方法
echo "🗑️  删除所有表..."
mysql -h"$DB_HOST" -P"$DB_PORT" -u"$DB_USER" -p"$DB_PASS" "$DB_NAME" -e "
SET FOREIGN_KEY_CHECKS = 0;
SHOW TABLES;
" | grep -v "Tables_in" | while read table; do
    if [ ! -z "$table" ]; then
        echo "删除表: $table"
        mysql -h"$DB_HOST" -P"$DB_PORT" -u"$DB_USER" -p"$DB_PASS" "$DB_NAME" -e "DROP TABLE IF EXISTS \`$table\`;"
    fi
done
mysql -h"$DB_HOST" -P"$DB_PORT" -u"$DB_USER" -p"$DB_PASS" "$DB_NAME" -e "SET FOREIGN_KEY_CHECKS = 1;"

if [ $? -eq 0 ]; then
    echo "✅ 所有表已删除"
else
    echo "❌ 删除表失败"
    exit 1
fi

echo ""

# 执行数据库迁移
echo "🔄 执行数据库迁移..."
go run main.go migrate -c config/settings.yml
if [ $? -eq 0 ]; then
    echo "✅ 数据库迁移完成"
else
    echo "❌ 数据库迁移失败"
    exit 1
fi

echo ""

# 清空所有表数据（避免重复键错误）
echo "🧹 清空表数据..."
mysql -h"$DB_HOST" -P"$DB_PORT" -u"$DB_USER" -p"$DB_PASS" "$DB_NAME" -e "
SET FOREIGN_KEY_CHECKS = 0;
TRUNCATE TABLE sys_user;
TRUNCATE TABLE sys_role;
TRUNCATE TABLE sys_menu;
TRUNCATE TABLE sys_role_menu;
TRUNCATE TABLE sys_dept;
TRUNCATE TABLE sys_post;
TRUNCATE TABLE sys_buildings;
TRUNCATE TABLE sys_house_types;
TRUNCATE TABLE sys_business_areas;
TRUNCATE TABLE sys_districts;
SET FOREIGN_KEY_CHECKS = 1;
"
if [ $? -eq 0 ]; then
    echo "✅ 表数据已清空"
else
    echo "❌ 清空表数据失败"
    exit 1
fi

echo ""

# 执行初始化文件
echo "📋 执行初始化文件..."
mysql -h"$DB_HOST" -P"$DB_PORT" -u"$DB_USER" -p"$DB_PASS" "$DB_NAME" < config/sql/init/01-mysql-begin.sql
if [ $? -eq 0 ]; then
    echo "✅ 执行: 01-mysql-begin.sql"
else
    echo "❌ 执行失败: 01-mysql-begin.sql"
    exit 1
fi

# 执行数据文件（按依赖顺序）
echo ""
echo "📊 执行数据文件..."
# 1. 先执行基础数据（没有外键依赖）
for sql_file in 01-buildings-data.sql sys_dept.sql sys_post.sql sys_role.sql sys_user.sql; do
    if [ -f "config/sql/data/$sql_file" ]; then
        echo "执行: $sql_file"
        if mysql -h"$DB_HOST" -P"$DB_PORT" -u"$DB_USER" -p"$DB_PASS" "$DB_NAME" < "config/sql/data/$sql_file"; then
            echo "✅ 成功: $sql_file"
        else
            echo "❌ 失败: $sql_file"
            exit 1
        fi
    else
        echo "⚠️  文件不存在: $sql_file"
    fi
done

# 2. 再执行菜单数据（依赖角色表）
if [ -f "config/sql/data/sys_menu.sql" ]; then
    echo "执行: sys_menu.sql"
    if mysql -h"$DB_HOST" -P"$DB_PORT" -u"$DB_USER" -p"$DB_PASS" "$DB_NAME" < "config/sql/data/sys_menu.sql"; then
        echo "✅ 成功: sys_menu.sql"
    else
        echo "❌ 失败: sys_menu.sql"
        exit 1
    fi
else
    echo "⚠️  文件不存在: sys_menu.sql"
fi

# 3. 跳过sys_role_menu.sql（因为sys_menu.sql中已经包含了角色菜单关联数据）
echo "⚠️  跳过: sys_role_menu.sql (角色菜单关联已在sys_menu.sql中处理)"

# 执行结束文件
echo ""
echo "📋 执行结束文件..."
mysql -h"$DB_HOST" -P"$DB_PORT" -u"$DB_USER" -p"$DB_PASS" "$DB_NAME" < config/sql/init/02-mysql-end.sql
if [ $? -eq 0 ]; then
    echo "✅ 执行: 02-mysql-end.sql"
else
    echo "❌ 执行失败: 02-mysql-end.sql"
    exit 1
fi

echo ""
echo "🎉 数据库重置完成!"
echo "结束时间: $(date)"
echo ""
echo "📝 数据库状态:"
echo "✅ 所有旧表已删除"
echo "✅ 新表结构已创建"
echo "✅ 基础数据已插入"
echo ""
echo "🚀 下一步操作:"
echo "1. 启动API服务器: go run main.go api -c config/settings.yml -p 8002"
echo "2. 测试API接口"
echo "3. 检查数据是否正确插入"
