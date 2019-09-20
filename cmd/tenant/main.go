package main

import (
	"github.com/decentralized-cloud/tenant/internal/cmd"
	"github.com/decentralized-cloud/tenant/pkg/util"
)

func main() {
	rootCmd := cmd.NewRootCommand()
	err := rootCmd.Execute()

	if err != nil {
		util.PrintError(err.Error())
	}
}
