package settings

import (
	"github.com/mannie-exe/packwiz-tx/cmd"
	"github.com/spf13/cobra"
)

// settingsCmd represents the base command when called without any subcommands
var settingsCmd = &cobra.Command{
	Use:   "settings",
	Short: "Manage pack settings",
}

func init() {
	cmd.Add(settingsCmd)
}
