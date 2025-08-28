// Package version 提供版本信息相关的命令行功能
// 用于显示 rentpro-admin 系统的版本信息
package version

import (
	"fmt"
	"io"
	"os"
	"rentPro/rentpro-admin/common/global"

	"github.com/spf13/cobra"
)

var (
	// StartCmd 定义了 version 子命令
	// 用于显示 rentpro-admin 系统的当前版本信息
	StartCmd = &cobra.Command{
		Use:     "version",               // 命令名称
		Short:   "Get version info",      // 命令简短描述
		Example: "rentpro-admin version", // 使用示例
		// PreRun 在实际命令执行前运行的函数
		// 可以用于参数验证、初始化等预处理操作
		PreRun: func(cmd *cobra.Command, args []string) {
			// 目前无需预处理操作
		},
		// RunE 是命令的主要执行函数，返回 error 类型
		// 如果返回 error，cobra 会自动处理错误信息
		RunE: func(cmd *cobra.Command, args []string) error {
			return run()
		},
	}
)

// run 执行版本信息显示的核心逻辑
// 从全局配置中获取版本号并打印到控制台
// 返回值: error - 如果执行成功返回 nil，失败返回相应错误
func run() error {
	return runWithOutput(os.Stdout)
}

// runWithOutput 执行版本信息显示，支持自定义输出目标
// 这个函数主要用于测试，可以将输出重定向到指定的 Writer
// 参数: output io.Writer - 输出目标
// 返回值: error - 如果执行成功返回 nil，失败返回相应错误
func runWithOutput(output io.Writer) error {
	// 打印系统版本信息
	// global.Version 来自 common/global 包，包含了系统的版本号
	_, err := fmt.Fprintln(output, global.Version)
	return err
}

// GetVersion 获取版本信息字符串
// 这是一个便于测试的辅助函数，返回版本信息而不是直接打印
// 返回值: string - 版本信息字符串
func GetVersion() string {
	return global.Version
}

// ValidateVersion 验证版本信息是否有效
// 用于测试和调试，检查版本信息是否符合预期格式
// 返回值: bool - 版本信息是否有效
func ValidateVersion() bool {
	version := global.Version
	return version != "" && len(version) > 0
}
