// Package cmd implements different commands that can be executed against tenant service
package cmd

import (
	"github.com/spf13/cobra"

	"github.com/decentralized-cloud/tenant/pkg/util"
)

func NewRootCommand() *cobra.Command {
	var cmd = &cobra.Command{
		Use: "tenant",
		PreRun: func(cmd *cobra.Command, args []string) {
			printHeader()
		},
	}

	// Register all commands
	cmd.AddCommand(
		newStartCommand(),
		newVersionCommand(),
	)

	return cmd
}

func printHeader() {
	util.PrintInfo("Tenant Serice")
}
