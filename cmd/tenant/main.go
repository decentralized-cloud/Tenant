package main

import (
	"github.com/decentralized-cloud/Tenant/internal/cmd"
	"github.com/decentralized-cloud/Tenant/pkg/util"
)

func main() {
	rootCmd := cmd.NewRootCommand()
	err := rootCmd.Execute()

	if err != nil {
		util.PrintError(err.Error())
	}
}
