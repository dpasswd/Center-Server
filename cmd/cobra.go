package cmd

import (
	"dh-passwd/cmd/api"
	"dh-passwd/cmd/config"
	"dh-passwd/cmd/migrate"
	"dh-passwd/cmd/version"
	"dh-passwd/tools"
	"errors"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// 自定义命令，定义run命令，调用run方法
var rootCmd = &cobra.Command{
	Use:   "dh-passwd",
	Short: "这是一个测试项目",
	Long:  `这是一个测试项目，旨在帮助大家管理密码`,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			tip()
			return errors.New(tools.Red("requires at least one arg"))
		}
		return nil
	},
	PersistentPreRunE: func(*cobra.Command, []string) error { return nil },
	Run: func(cmd *cobra.Command, args []string) {
		tip()
	},
}

func tip() {
	usageStr := `欢迎进入 ` + tools.Green(`这个测试项目`) + ` 可以使用 ` + tools.Red(`-h`) + ` 查看命令`
	fmt.Printf("%s\n", usageStr)
}

func init() {
	rootCmd.AddCommand(api.StartCmd)
	rootCmd.AddCommand(config.StartCmd)
	rootCmd.AddCommand(migrate.StartCmd)
	rootCmd.AddCommand(version.StartCmd)
}

//Execute : apply commands
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(-1)
	}
}
