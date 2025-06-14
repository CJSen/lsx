package cmd

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/CJSen/lsx/config"
	"github.com/CJSen/lsx/utils"
	markdown "github.com/MichaelMure/go-term-markdown"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

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
	ErrCommandNotFound = errors.New("command not found")
)

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
			fmt.Printf("[sorry]: failed to download command <%s>\n", cmd)
		}
	}

	p := filepath.Join(cfg.DataDir, fmt.Sprintf("%s.md", cmd))
	if !utils.FileExists(p) {
		err := utils.RetryDownloadFile(url, path, cmd)
		if err != nil {
			fmt.Printf("[sorry]: failed to retrieve command <%s>\n", cmd)
			return
		}
		if errors.Is(err, ErrCommandNotFound) {
			fmt.Printf("[sorry]: could not found command <%s>\n", cmd)
			return
		}
	}

	source, err := os.ReadFile(p)
	if err != nil {
		fmt.Printf("[sorry]: failed to open file <%s>\n", p)
		return
	}
	markdown.BlueBgItalic = color.New(color.FgBlue).SprintFunc()
	result := markdown.Render(string(source), 80, 6)
	fmt.Println(string(result))
}
