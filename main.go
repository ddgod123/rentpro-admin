package main

import (
	"fmt"

	"rentPro/rentpro-admin/cmd"
)

func main() {
	fmt.Println("hello world")
	cmd.Execute()
}

/*

同步仓库：
./sync_repositories.sh --dry-run

启动API服务器：go run main.go api -p 8002
• 执行数据库迁移：go run main.go migrate
• 查看配置信息：go run main.go config


# 启动API服务器
go run main.go api -c config/settings.yml -p 8002

# 执行数据库迁go run main.go migrate -c config/settings.yml

# 查看配置信息
go run main.go config -c config/settings.yml

# 查看版本信息
go run main.go version



*/
