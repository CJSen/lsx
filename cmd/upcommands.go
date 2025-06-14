package cmd

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

const commandDataUrl = commandUrlBase + "/dist/data.json"

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

func getPath() (string, error) {
	c, err := getConfigContent()
	if err != nil {
		return "", err
	}

	if err := makeCmdDir(c.Dir); err != nil {
		return "", err
	}
	path := filepath.Join(c.Dir, linuxCommandJson)
	return path, nil
}

func LoadCommandJson(path string) ([]byte, error) {
	if isFileExist(path) {
		data, err := os.ReadFile(path)
		if err != nil {
			return nil, fmt.Errorf("failed to read existing linux-command.json: %w", err)
		}
		return data, nil
	} else {
		err := os.WriteFile(path, linuxCommandJsonTemp, 0666)
		if err != nil {
			return nil, err
		}
		return linuxCommandJsonTemp, nil
	}
}

func CheckCommandJson() ([]byte, error) {

	linuxCommandJsonPath, err := getPath()
	if err != nil {
		return nil, fmt.Errorf("failed to get linux-command.json path: %w", err)
	}
	data, err := LoadCommandJson(linuxCommandJsonPath)
	if err != nil {
		return nil, fmt.Errorf("failed to load linux-command.json: %w", err)
	}
	return data, nil
}

func downloadLatestJSON() error {
	resp, err := http.Get(commandDataUrl)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	path, err := getPath()
	if err != nil {
		return fmt.Errorf("failed to get path for linux-command.json: %w", err)
	}

	outFile, err := os.Create(path)
	if err != nil {
		return err
	}
	defer outFile.Close()

	_, err = io.Copy(outFile, resp.Body)
	if err != nil {
		return fmt.Errorf("failed to write to linux-command.json: %w", err)

	} else {
		LoadCommandJson(path)
	}
	return err
}
