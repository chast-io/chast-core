package cmd

import (
	"chast.io/core/pkg/api/refactoring"
	util "chast.io/core/pkg/util/fs/file"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// refactoringCmd represents the refactoring command.
var refactoringCmd = &cobra.Command{ //nolint:exhaustruct // Only defining required fields
	Use:   "refactoring <chastConfigFile>",
	Short: "A brief description of your command", // TODO
	Long:  ``,                                    // TODO
	Args:  cobra.MatchAll(cobra.MinimumNArgs(1)),
	Run: func(cmd *cobra.Command, args []string) {
		recipeFileArg := args[0]

		file, newFileError := util.NewFile(recipeFileArg)
		if newFileError != nil || !file.Exists() {
			log.Fatalf("Recipe file \"%v\" does not exist.\n", file.AbsolutePath)
		}
		refactoring.Run(file, args[1:]...)
	},
}

func init() { //nolint:gochecknoinits // This is the way cobra wants it.
	runCmd.AddCommand(refactoringCmd)

	defaultHelpFunction := refactoringCmd.HelpFunc()
	refactoringCmd.SetHelpFunc(func(cmd *cobra.Command, args []string) { helpFunction(cmd, args, defaultHelpFunction) })

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:

	// Cobra supports local flags which will only refactoring when this command
	// is called directly, e.g.:
	// refactoringCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func helpFunction(cmd *cobra.Command, args []string, defaultHelpFunction func(*cobra.Command, []string)) {
	if len(cmd.ValidArgs) > 0 {
		helpFunction(cmd, args, defaultHelpFunction)
	} else {
		defaultHelpFunction(cmd, args)
	}
}
