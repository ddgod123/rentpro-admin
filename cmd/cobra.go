package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

// rootCmd 定义了应用程序的根命令
var rootCmd = &cobra.Command{
	Use:          "rentpro-admin",
	Short:        "Cobra is a CLI library ",
	SilenceUsage: true,
	Long:         `Cobra is a CLI library for Go that empowers applications.`,

	// Args: func(cmd *cobra.Command, args []string) error {
	// 	if len(args) < 1 {
	// 		return errors.New("requires at least one arg")
	// 	}
	// 	return nil
	// },

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

func init() {

	// rootCmd.AddCommand(api.StartCmd)
	// rootCmd.AddCommand(migrate.StartCmd)
	// rootCmd.AddCommand(version.StartCmd)
	// rootCmd.AddCommand(config.StartCmd)
	// rootCmd.AddCommand(app.StartCmd)

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
