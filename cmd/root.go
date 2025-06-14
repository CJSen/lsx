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

func init() {
	initData() // 初始化数据和配置
	cmdMap := make(map[string]interface{})
	if err := json.Unmarshal(linuxCommandJSON, &cmdMap); err != nil {
		panic(fmt.Sprintf("failed to parse linux-command.json: %v\ncontent: %.128s", err, string(linuxCommandJSON)))
	}
	commands = make([]string, 0, len(cmdMap))
	for k := range cmdMap {
		commands = append(commands, k)
	}

	// 注册各子命令
	rootCmd.AddCommand(
		NewShowCommand(),
		NewUpgradeCommand(),
		NewVersionCommand(),
		NewSearchCommand(),
		NewUpdateCommand(),
	)
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
