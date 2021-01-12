// Package cmd implements different commands that can be executed against tenant service
package cmd

import (
	"fmt"
	"time"

	"github.com/decentralized-cloud/tenant/pkg/util"
	gocoreUtil "github.com/micro-business/go-core/pkg/util"
	"github.com/spf13/cobra"
)

func newStartCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "start",
		Short: "Start the Tenant service",
		Run: func(cmd *cobra.Command, args []string) {
			gocoreUtil.PrintInfo(fmt.Sprintf("Copyright (C) %d, Micro Business Ltd.\n", time.Now().Year()))
			gocoreUtil.PrintYAML(gocoreUtil.GetVersion())
			util.StartService()
		},
	}
}
