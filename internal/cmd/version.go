// Package cmd implements different commands that can be executed against project service
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
		Short: "Get Project CLI version",
		Run: func(cmd *cobra.Command, args []string) {
			util.PrintInfo("Project CLI\n")
			util.PrintInfo(fmt.Sprintf("Copyright (C) %d, Micro Business Ltd.\n", time.Now().Year()))
			util.PrintYAML(util.GetVersion())
		},
	}
}
