package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
)

// commandIndex 结构体，映射命令索引数据
// 对应 https://raw.githubusercontent.com/jaywcjlove/linux-command/master/dist/data.json
// n: 命令名，p: 路径，d: 描述
type commandIndex struct {
	Name        string `json:"n"`
	Path        string `json:"p"`
	Description string `json:"d"`
}

// 创建 search 子命令
func NewSearchCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "search <command>",
		Short: "Search command by keywords",
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) == 0 {
				fmt.Println("[sorry]: the search command does not accept any keywords")
				return
			}
			searchCmd(args[0])
		},
	}
	return cmd
}

// 反序列化命令索引数据
func unmarshalIndex() (map[string]commandIndex, error) {
	ret := make(map[string]commandIndex)
	err := json.Unmarshal(linuxCommandJSON, &ret)
	return ret, err
}

// 执行命令搜索
func searchCmd(keywords string) {
	keywords = strings.ToLower(keywords)

	table := tablewriter.NewWriter(os.Stdout)
	table.Append([]string{"command", "description"})

	found := false
	index, err := unmarshalIndex()
	if err != nil {
		fmt.Println("[error]: failed to parse command index:", err)
		return
	}
	for _, v := range index {
		if strings.Contains(strings.ToLower(v.Name), keywords) || strings.Contains(strings.ToLower(v.Description), keywords) {
			table.Append([]string{v.Name, v.Description})
			found = true
		}
	}
	if found {
		table.Render()
	} else {
		fmt.Println("[sorry]: no command found for", keywords)
	}
}
