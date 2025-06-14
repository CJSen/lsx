package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var versionInfo = "dev" // 版本信息，默认 dev

// 创建 version 子命令
func NewVersionCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "version",
		Short: "Prints the version about lsx",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("lsx version: %s\n", versionInfo)
		},
	}
	return cmd
}
