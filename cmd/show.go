package cmd

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/CJSen/lsx/config"
	"github.com/CJSen/lsx/utils"
	markdown "github.com/MichaelMure/go-term-markdown"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

// 创建 show 子命令
func NewShowCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "show <command>",
		Short: "Show the specified command usage.",
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) == 0 {
				fmt.Println("[sorry]: the show command does not accept any arguments")
				return
			}
			force, _ := cmd.Flags().GetBool("force")
			showCmd(args[0], force)
		},
	}
	cmd.Flags().BoolP("force", "f", false, "force to refresh command usage from remote.")
	return cmd
}

var (
	ErrCommandNotFound = errors.New("command not found") // 未找到命令错误
)

// 展示指定命令的用法，支持强制刷新
func showCmd(cmd string, force bool) {
	cfg := config.GlobalConfig
	cmd = strings.ToLower(cmd)

	url := cfg.RemoteBaseUrl + fmt.Sprintf("/command/%s.md", cmd)
	path := filepath.Join(cfg.DataDir, fmt.Sprintf("%s.md", cmd))
	if force {
		if err := utils.RetryDownloadFile(url, path, cmd); err != nil {
			if errors.Is(err, ErrCommandNotFound) {
				fmt.Printf("[sorry]: could not found command <%s>\n", cmd)
				return
			}
			fmt.Printf("[sorry]: failed to download command <%s>: %v\n", cmd, err)
			return
		}
	}

	p := filepath.Join(cfg.DataDir, fmt.Sprintf("%s.md", cmd))
	if !utils.FileExists(p) {
		fmt.Printf("[sorry]: could not found command <%s>, it will be try to download \n", cmd)
		err := utils.RetryDownloadFile(url, path, cmd)
		if err != nil {
			fmt.Printf("[sorry]: failed to retrieve command <%s>: %v\n", cmd, err)
			return
		}
		if errors.Is(err, ErrCommandNotFound) {
			fmt.Printf("[sorry]: could not found command <%s>\n", cmd)
			return
		}
	}

	source, err := os.ReadFile(p)
	if err != nil {
		fmt.Printf("[sorry]: failed to open file <%s>: %v\n", p, err)
		return
	}
	result := ""
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("[error]: failed to render markdown.")
			fmt.Println(string(source))
		}
	}()
	markdown.BlueBgItalic = color.New(color.FgBlue).SprintFunc()
	result = string(markdown.Render(string(source), 80, 6))

	if config.GlobalConfig != nil && config.GlobalConfig.UseLess {
		cmd := exec.Command("less")
		cmd.Stdin = strings.NewReader(result)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		_ = cmd.Run()
	} else {
		fmt.Println(result)
	}
}
