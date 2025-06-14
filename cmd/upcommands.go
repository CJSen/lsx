package cmd

import (
	_ "embed"
	"fmt"
	"os"
	"path/filepath"

	"github.com/CJSen/lsx/config"
	"github.com/CJSen/lsx/utils"
	"github.com/spf13/cobra"
)

//go:embed linux-command.json
var linuxCommandJsonTemp []byte

func NewUpdateCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "upcommands",
		Short: "Update the embedded linux-command.json to the latest version.",
		Run: func(cmd *cobra.Command, args []string) {
			err := downloadLatestJSON()
			if err != nil {
				fmt.Printf("[error]: failed to update linux-command.json: %v\n", err)
			} else {
				fmt.Println("[success]: linux-command.json updated.")
			}
		},
	}
}

func CheckCommandJson() ([]byte, error) {
	cfg := config.GlobalConfig
	linuxCommandJsonPath := filepath.Join(cfg.DataDir, "linux-command.json")
	if utils.FileExists(linuxCommandJsonPath) {
		data, err := os.ReadFile(linuxCommandJsonPath)
		if err != nil {
			return nil, fmt.Errorf("failed to read existing linux-command.json: %w", err)
		}
		return data, nil
	} else {
		err := os.WriteFile(linuxCommandJsonPath, linuxCommandJsonTemp, 0666)
		if err != nil {
			return nil, err
		}
		return linuxCommandJsonTemp, nil
	}
}

func downloadLatestJSON() error {
	cfg := config.GlobalConfig
	commandDataUrl := cfg.RemoteBaseUrl + "/dist/data.json"
	linuxCommandJsonPath := filepath.Join(cfg.DataDir, "linux-command.json")
	err := utils.RetryDownloadFile(commandDataUrl, linuxCommandJsonPath, "linux-command.json")
	if err != nil {
		return err
	}
	return nil
}
