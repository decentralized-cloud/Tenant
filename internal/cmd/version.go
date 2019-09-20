// Package cmd implements different commands that can be executed against tenant service
package cmd

import (
	"github.com/decentralized-cloud/Tenant/pkg/util"
	"github.com/spf13/cobra"
)

func newVersionCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Get Tenant CLI version",
		Run: func(cmd *cobra.Command, args []string) {
			util.PrintInfo("Tenant CLI\n")
			util.PrintInfo("Copyright (C) 2019, Micro Business Ltd.\n")
			util.PrintYAML(util.GetVersion())
		},
	}
}
