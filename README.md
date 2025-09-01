# rentpro-admin



# 楼盘管理 
 -- 户型 、户型图 



git add .
git commit -m "Initial commit"
git push -u origin main


git branch -M main



hahahahah


git init   # 创建新的git仓库
git status # 查看状态
git branch # 查看分支
git branch dev  # 创建dev分支
git branch -d dev  # 删除 dev 分支
git branch -a # 查看远程分支
git checkout -b dev # 创建dev分支，并切换到dev分支
git checkout master  # 切换到master分支
git add filename  # 添加指定文件，把当前文件放入暂存区域
git add .  # 表示添加新文件和编辑过的文件不包括删除的文件
git add -A  # 表示添加所有内容
git commit  # 给暂存区域生成快照并提交
git reset -- files # 用来撤销最后一次 git add files，也可以用 git reset 撤销所有暂存区域文件
git push origin master  # 推送改动到master分支（前提是已经clone了现有仓库）
git remote add origin   # 没有克隆现有仓库，想仓库连接到某个远程服务器
git pull  # 更新本地仓库到最新版本（多人合作的项目），以在我们的工作目录中 获取（fetch） 并 合并（merge） 远端的改动
git diff    # 查看两个分支差异
git diff  # 查看已修改的工作文档但是尚未写入缓冲的改动
git rm   # 用于简单的从工作目录中手工删除文件
git rm -f   # 删除已经修改过的并且放入暂存区域的文件，必须使用强制删除选项 -f
git mv   # 用于移动或重命名一个文件、目录、软链接
git log  # 列出历史提交记录




go get -u github.com/gin-gonic/gin
go get -u github.com/spf13/viper
go get -u github.com/appleboy/gin-jwt/v2
go get -u github.com/google/uuid

go get -u github.com/casbin/casbin/v2
go get github.com/casbin/casbin/v2
go get -u github.com/swaggo/swag/cmd/swag
go get -u github.com/swaggo/gin-swagger
go get -u github.com/swaggo/files




# ok
go get -u gorm.io/gorm
go get -u gorm.io/driver/mysql

//要获取 Cobra 命令行库，可以使用以下 go get 命令：
go get -u github.com/spf13/cobra
go get -u github.com/spf13/cobra-cli
go get -u github.com/spf13/pflag
go get -u github.com/spf13/viper

# 创建 app 目录及其子目录
mkdir -p app/admin
mkdir -p app/jobs
mkdir -p app/other

# 创建 cmd 目录及其子目录
mkdir -p cmd/api
mkdir -p cmd/app
mkdir -p cmd/config
mkdir -p cmd/migrate
mkdir -p cmd/version

# 创建 common 目录及其子目录
mkdir -p common/actions
mkdir -p common/apis
mkdir -p common/database
mkdir -p common/dto
mkdir -p common/file_store
mkdir -p common/global
mkdir -p common/middleware
mkdir -p common/models
mkdir -p common/response
mkdir -p common/service
mkdir -p common/storage

# 创建其他必要目录
mkdir -p config
mkdir -p docs
mkdir -p static
mkdir -p template
mkdir -p test

# 确认目录结构创建完成
ls -la



完整调用链 ：
- 命令行执行 go run main.go migrate -c config/settings.yml 启动迁移
- cmd/migrate/server.go 中的 run() 函数调用 initDB()
- initDB() 函数调用 migrateModel()
- migrateModel() 函数设置数据库连接并调用 migration.Migrate.Migrate()
- migration.Migrate.Migrate() 按顺序执行所有注册的迁移函数
- 在 _1599190683659Tables 迁移函数中调用 models.InitDb(tx)





