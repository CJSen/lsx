package cmd

import (
	_ "embed"
	"encoding/json"
	"fmt"

	"github.com/CJSen/lsx/config"
	"github.com/CJSen/lsx/utils"
	"github.com/spf13/cobra"
)

// commands 保存所有可用命令名
var commands []string

// linuxCommandJSON 保存命令数据内容
var linuxCommandJSON []byte

// 根命令，lsx 的主入口
var rootCmd = &cobra.Command{
	Use:   "lsx",
	Short: "Impressive Linux commands cheat sheet cli.",
}

// Execute 启动命令行程序
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
	}
}

// initData 初始化配置和命令数据
func initData() {
	// 初始化配置
	_ = config.ParseConfig()

	// 确保数据目录存在，否则会报错
	err := utils.MakesureDir(config.GlobalConfig.DataDir)
	if err != nil {
		panic(fmt.Sprintf("failed to create dir: %v", err))
	}
	var err2 error
	// 检查并加载命令数据
	linuxCommandJSON, err2 = CheckCommandJson()
	if err2 != nil {
		panic(fmt.Sprintf("failed to parse linux-command.json: %v", err2))
	} else if linuxCommandJSON == nil {
		panic("failed to parse linux-command.json: file not found")
	}
}

func init() {
	// 初始化数据和配置
	initData()
	// 自定义 help 输出
	rootCmd.SetHelpFunc(func(cmd *cobra.Command, args []string) {
		fmt.Println("Impressive Linux commands cheat sheet cli.")
		fmt.Println("Usage:")
		fmt.Println("  lsx [command]")
		fmt.Println("\nAvailable Commands:")
		if !config.GlobalConfig.UseShow {
			fmt.Println("  your-cmd    Enter the specified command to show its usage.")
		}
		if config.GlobalConfig.UseShow {
			fmt.Println("  show        Show the specified command usage.")
		}
		fmt.Println("  search      Search command by keywords")
		fmt.Println("  upcommands  Update the embedded linux-command.json to the latest version.")
		fmt.Println("  upgrade     Upgrade all commands from remote.")
		fmt.Println("  version     Prints the version about lsx")
		fmt.Println("  help        Help about any command")
	})
	cmdMap := make(map[string]interface{})
	if err := json.Unmarshal(linuxCommandJSON, &cmdMap); err != nil {
		panic(fmt.Sprintf("failed to parse linux-command.json: %v\ncontent: %.128s", err, string(linuxCommandJSON)))
	}
	commands = make([]string, 0, len(cmdMap))
	for k := range cmdMap {
		commands = append(commands, k)
	}

	// 根据配置决定注册方式
	if config.GlobalConfig.UseShow {
		// 兼容旧方式，注册 show 子命令
		rootCmd.AddCommand(
			NewShowCommand(),
			NewUpgradeCommand(),
			NewVersionCommand(),
			NewSearchCommand(),
			NewUpdateCommand(),
		)
	} else {
		// 注册所有一级命令
		for _, cmdName := range commands {
			cmd := &cobra.Command{
				Use:   cmdName,
				Short: fmt.Sprintf("Show usage for %s", cmdName),
				Run: func(cmd *cobra.Command, args []string) {
					showCmd(cmd.Use, false)
				},
			}
			rootCmd.AddCommand(cmd)
		}
		// 其他功能性命令依然注册
		rootCmd.AddCommand(
			NewUpgradeCommand(),
			NewVersionCommand(),
			NewSearchCommand(),
			NewUpdateCommand(),
		)
	}
}
