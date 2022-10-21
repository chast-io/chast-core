package cmd

import (
	"github.com/spf13/cobra"
)

import (
	"chast.io/core/pkg/api/refactoring"
	util "chast.io/core/pkg/util"
)

// refactoringCmd represents the refactoring command
var refactoringCmd = &cobra.Command{
	Use:   "refactoring",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		//println(fmt.Sprintf("%#v", cmd.Flags()))
		//println(fmt.Sprintf("%#v", args))

		file := util.NewFile(args[0])
		refactoring.Run(file)
	},
}

func init() {
	runCmd.AddCommand(refactoringCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// refactoringCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only refactoring when this command
	// is called directly, e.g.:
	// refactoringCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
