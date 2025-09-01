package cmd

import (
	"errors"
	"fmt"
	"os"

	"rentPro/rentpro-admin/cmd/api"
	"rentPro/rentpro-admin/cmd/config"
	"rentPro/rentpro-admin/cmd/migrate"
	"rentPro/rentpro-admin/cmd/version"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

// rootCmd 定义了应用程序的根命令
var rootCmd = &cobra.Command{
	Use:          "rentpro-admin",
	Short:        "Cobra is a CLI library ",
	SilenceUsage: true,
	Long:         `Cobra is a CLI library for Go that empowers applications.`,

	// Args 函数用于验证命令行参数
	// 要求用户必须提供至少一个子命令（如 version、config 等）
	// 这确保了用户不能仅运行 "rentpro-admin" 而不指定具体操作
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("请指定一个子命令，如: version, config")
		}
		return nil
	},

	PersistentPreRunE: func(*cobra.Command, []string) error { return nil },

	Run: func(cmd *cobra.Command, args []string) {
		// 打印命令信息
		fmt.Println("=== 命令执行信息 ===")
		fmt.Printf("命令名称: %s\n", cmd.Name())
		fmt.Printf("完整命令路径: %s\n", cmd.CommandPath())
		fmt.Printf("输入的参数: %v\n", args)
		fmt.Printf("参数数量: %d\n", len(args))

		// 打印所有设置的标志(flags)
		hasFlags := false
		cmd.Flags().VisitAll(func(flag *pflag.Flag) {
			if flag.Changed {
				if !hasFlags {
					fmt.Println("\n设置的标志:")
					hasFlags = true
				}
				fmt.Printf("  --%s = %s\n", flag.Name, flag.Value)
			}
		})

		if !hasFlags {
			fmt.Println("\n没有设置任何标志")
		}

		// 打印原始命令行
		fmt.Printf("\n原始命令行参数: %v\n", os.Args)
		fmt.Println("==================")
	},
}

// func init 用来注册命令
func init() {

	// 注册 migrate 子命令到根命令
	// migrate.StartCmd 来自 cmd/migrate/server.go，提供数据库迁移功能
	// 注册后用户可以通过 "rentpro-admin migrate" 来执行数据库迁移操作
	// 支持的用法: rentpro-admin migrate -c config/settings.yml
	rootCmd.AddCommand(migrate.StartCmd)

	// 注册 version 子命令到根命令
	// version.StartCmd 来自 cmd/version/server.go，提供版本信息显示功能
	// 注册后用户可以通过以下方式查看版本信息：
	//   - rentpro-admin version           : 显示详细版本信息
	//   - rentpro-admin migrate -v        : 在migrate命令中显示版本
	//   - rentpro-admin migrate --version : 在migrate命令中显示版本
	// 版本信息来源: common/global/adm.go 中的 Version 常量 (当前: "2.2.0")
	rootCmd.AddCommand(version.StartCmd)

	// 注册 config 子命令到根命令
	// config.StartCmd 来自 cmd/config/server.go，提供配置信息显示功能
	// 注册后用户可以通过以下方式查看和验证配置：
	//   - rentpro-admin config -c config/settings.yml : 显示指定配置文件的内容
	//   - rentpro-admin config -v                     : 显示config工具版本
	//   - rentpro-admin config --version               : 显示config工具版本
	// 功能特性: 支持YAML配置解析、敏感信息隐藏、配置验证等
	// 配置来源: 读取项目的 settings.yml 配置文件
	rootCmd.AddCommand(config.StartCmd)

	// 注册 api 子命令到根命令
	// api.StartCmd 来自 cmd/api/server.go，提供HTTP API服务器功能
	// 注册后用户可以通过以下方式启动API服务器：
	//   - rentpro-admin api -c config/settings.yml : 使用指定配置文件启动API服务器
	//   - rentpro-admin api                       : 使用默认配置启动API服务器
	// 功能特性: 提供完整的权限管理API、用户认证、JWT令牌等
	// 服务端口: 默认8000（可在config/settings.yml中配置）
	rootCmd.AddCommand(api.StartCmd)

}

// Execute 是命令行应用的入口函数，由main.go调用
// 它执行根命令并处理可能出现的错误
func Execute() {

	fmt.Println("main -----> Execute()")

	// 在执行命令前打印原始输入
	fmt.Printf("接收到的命令行参数---------------- %v\n", os.Args)

	if err := rootCmd.Execute(); err != nil {

		os.Exit(-1)
	}

}
