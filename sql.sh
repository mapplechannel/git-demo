#!/bin/bash

# 数据库连接信息
DB_HOST="172.21.44.157"
DB_PORT="7654"
DB_USER="postgres"
DB_PASS="postgres"  # 直接在脚本中设置密码
DB_NAME="postgres"
DB_SCHEMA="hsm_scheduling"
TABLE_NAME="tasks"
OUTPUT_DIR="/path/to/output"  # 请根据需要修改文件保存路径
BACKUP_FILE="$OUTPUT_DIR/backup_tasks.sql"

# 1. 使用 pg_dump 备份数据库
echo "开始备份数据库：$DB_NAME 中的表：$TABLE_NAME"

# 直接在脚本中设置 PGPASSWORD
export PGPASSWORD=$DB_PASS

# 使用 pg_dump 进行备份
pg_dump -h $DB_HOST -p $DB_PORT -U $DB_USER -d $DB_NAME -n $DB_SCHEMA -t $DB_SCHEMA.$TABLE_NAME --no-password > $BACKUP_FILE

if [ $? -ne 0 ]; then
    echo "数据库备份失败!"
    exit 1
fi

echo "数据库备份完成，备份文件：$BACKUP_FILE"

# 2. 使用 psql 执行删除操作
echo "开始删除 integrated 字段中以 ioit 开头的数据"
DELETE_QUERY="DELETE FROM $DB_SCHEMA.$TABLE_NAME WHERE integrated LIKE 'ioit%'"

# 使用 psql 执行 SQL 删除语句
psql -h $DB_HOST -p $DB_PORT -U $DB_USER -d $DB_NAME -c "$DELETE_QUERY"

if [ $? -ne 0 ]; then
    echo "删除数据失败!"
    exit 1
fi

echo "删除操作完成，集成字段中以 ioit 开头的数据已删除。"

# 清除环境变量
unset PGPASSWORD

exit 0
