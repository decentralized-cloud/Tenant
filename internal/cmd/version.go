// Package cmd implements different commands that can be executed against tenant service
package cmd

import (
	"fmt"
	"time"

	"github.com/micro-business/go-core/pkg/util"
	"github.com/spf13/cobra"
)

func newVersionCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Get Tenant CLI version",
		Run: func(cmd *cobra.Command, args []string) {
			util.PrintInfo("Tenant CLI\n")
			util.PrintInfo(fmt.Sprintf("Copyright (C) %d, Micro Business Ltd.\n", time.Now().Year()))
			util.PrintYAML(util.GetVersion())
		},
	}
}
