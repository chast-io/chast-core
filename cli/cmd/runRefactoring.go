package cmd

import (
	"chast.io/core/pkg/api/refactoring"
	util "chast.io/core/pkg/util/fs/file"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// runRefactoringCmd represents the refactoring command.
var runRefactoringCmd = &cobra.Command{ //nolint:exhaustruct // Only defining required fields
	Use:   "refactoring <chastConfigFile>",
	Short: "Run a refactoring recipe",
	Long: `Run a refactoring recipe. 
The available flags and parameters are available by calling it with the --help flag [TODO].`,
	Args: cobra.MatchAll(cobra.MinimumNArgs(1)),
	Run: func(cmd *cobra.Command, args []string) {
		recipeFileArg := args[0]

		file, newFileError := util.NewFile(recipeFileArg)
		if newFileError != nil || !file.Exists() {
			log.Fatalf("Recipe file \"%v\" does not exist.", file.AbsolutePath)
		}
		refactoring.Run(file, args[1:]...)
	},
}

func init() { //nolint:gochecknoinits // This is the way cobra wants it.
	runCmd.AddCommand(runRefactoringCmd)

	defaultHelpFunction := runRefactoringCmd.HelpFunc()
	runRefactoringCmd.SetHelpFunc(func(cmd *cobra.Command, args []string) { runRefactoringHelpFunction(cmd, args, defaultHelpFunction) })
}

func runRefactoringHelpFunction(cmd *cobra.Command, args []string, defaultHelpFunction func(*cobra.Command, []string)) {
	if len(cmd.ValidArgs) > 0 {
		runRefactoringHelpFunction(cmd, args, defaultHelpFunction)
	} else {
		defaultHelpFunction(cmd, args)
	}
}
