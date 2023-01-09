package cmd

import (
	"github.com/spf13/cobra"
)

// runCmd represents the refactoring command.
var runCmd = &cobra.Command{ //nolint:exhaustruct // Only defining required fields
	Use:   "run",
	Short: "Run a certain type of recipe",
	//Long: ``,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.HelpFunc()(cmd, args)
	},
}

func init() { //nolint:gochecknoinits // This is the way cobra wants it.
	rootCmd.AddCommand(runCmd)
}
