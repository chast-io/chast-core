package cmd

import (
	"chast.io/core/pkg/api/refactoring"
	util "chast.io/core/pkg/util/fs/file"
	log "github.com/sirupsen/logrus"

	"github.com/spf13/cobra"
)

// testRefactoringCmd represents the test command
var testRefactoringCmd = &cobra.Command{
	Use:   "refactoring <chastConfigFile>",
	Short: "Test a refactoring recipe",
	Long: `This command tests a refactoring recipe based on the test section in the recipe itself.
The parameters and "input" files are passed as arguments and flags respectively.
The output is then compared against the expected output files in the "expected" folder.`,
	Args: cobra.MatchAll(cobra.MinimumNArgs(1)),
	Run: func(cmd *cobra.Command, args []string) {
		recipeFileArg := args[0]

		file, newFileError := util.NewFile(recipeFileArg)
		if newFileError != nil || !file.Exists() {
			log.Fatalf("Recipe file \"%v\" does not exist.", file.AbsolutePath)
		}
		refactoring.Test(file, args[1:]...)
	},
}

func init() { //nolint:gochecknoinits // This is the way cobra wants it.
	testCmd.AddCommand(testRefactoringCmd)

	defaultHelpFunction := testRefactoringCmd.HelpFunc()
	testRefactoringCmd.SetHelpFunc(func(cmd *cobra.Command, args []string) { testRefactoringHelpFunction(cmd, args, defaultHelpFunction) })
}

func testRefactoringHelpFunction(cmd *cobra.Command, args []string, defaultHelpFunction func(*cobra.Command, []string)) {
	if len(cmd.ValidArgs) > 0 {
		testRefactoringHelpFunction(cmd, args, defaultHelpFunction)
	} else {
		defaultHelpFunction(cmd, args)
	}
}
