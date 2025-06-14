package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var versionInfo = "dev"

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
