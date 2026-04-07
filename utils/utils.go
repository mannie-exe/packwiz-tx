package utils

import (
	"github.com/mannie-exe/packwiz-tx/cmd"
	"github.com/spf13/cobra"
)

// utilsCmd represents the base command when called without any subcommands
var utilsCmd = &cobra.Command{
	Use:   "utils",
	Short: "Utilities for managing packwiz itself",
}

func init() {
	cmd.Add(utilsCmd)
}
