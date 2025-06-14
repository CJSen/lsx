package cmd

import (
	_ "embed"
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"
)

//go:embed linux-command.json
var linuxCommandJsonTemp []byte

//go:embed version
var versionTemp string

var linuxCommandJSON []byte

const commandUrlBase = "https://unpkg.com/linux-command@latest"

var commands []string

var rootCmd = &cobra.Command{
	Use:   "lsx",
	Short: "Impressive Linux commands cheat sheet cli.",
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
	}
}

func init() {
	linuxCommandJSON, err := CheckCommandJson()
	if err != nil {
		panic("failed to parse linux-command.json: " + err.Error())
	} else if linuxCommandJSON == nil {
		panic("failed to parse linux-command.json: file not found")
	}

	var cmdMap map[string]interface{}
	if err := json.Unmarshal(linuxCommandJSON, &cmdMap); err != nil {
		panic("failed to parse linux-command.json: " + err.Error())
	}
	commands = make([]string, 0, len(cmdMap))
	for k := range cmdMap {
		commands = append(commands, k)
	}

	rootCmd.AddCommand(
		NewShowCommand(),
		NewUpgradeCommand(),
		NewVersionCommand(),
		NewSearchCommand(),
		NewUpdateCommand(),
	)
}
