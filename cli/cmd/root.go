package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string // config file currently not used

// rootCmd represents the base command when called without any subcommands.
var rootCmd = &cobra.Command{ //nolint:exhaustruct // Only defining required fields
	Use:   "chast",
	Short: "CHAnge STuff - A CLI for unifying tools and automating changes",
	Long: `
 $$$$$$\  $$\   $$\  $$$$$$\   $$$$$$\ $$$$$$$$\ 
$$  __$$\ $$ |  $$ |$$  __$$\ $$  __$$\\__$$  __|
$$ /  \__|$$ |  $$ |$$ /  $$ |$$ /  \__|  $$ |   
$$ |      $$$$$$$$ |$$$$$$$$ |\$$$$$$\    $$ |   
$$ |      $$  __$$ |$$  __$$ | \____$$\   $$ |   
$$ |  $$\ $$ |  $$ |$$ |  $$ |$$\   $$ |  $$ |   
\$$$$$$  |$$ |  $$ |$$ |  $$ |\$$$$$$  |  $$ |   
 \______/ \__|  \__|\__|  \__| \______/   \__|

Run refactorings and other commands through an unified system no matter which operating system, installer or programming language.

Required tools: 
- unionfs-fuse (Linux only, for Apple see MacOS support section in their README)
- user namespace support required
- (For OverlayFs-MergerFs-Isolation-Strategy: OverlayFs, Fuse, MergerFs required)`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() { //nolint:gochecknoinits // This is the way cobra wants it.
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.cmd.yaml)")

	// Cobra also supports local flags, which will only refactoring
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		// Search config in home directory with name ".cmd" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigType("yaml")
		viper.SetConfigName(".cmd")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		_, _ = fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	}
}
