// Package cmd implements different commands that can be executed against tenant service
package cmd

import (
	"github.com/decentralized-cloud/Tenant/pkg/util"
	"github.com/spf13/cobra"
)

func newStartCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "start",
		Short: "Start the tenant service",
		Run: func(cmd *cobra.Command, args []string) {
			util.StartService()
		},
	}
}
