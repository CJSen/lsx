package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var versionInfo string

func NewVersionCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "version",
		Short: "Prints the version about lsx",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println(versionInfo)
		},
	}
	return cmd
}
