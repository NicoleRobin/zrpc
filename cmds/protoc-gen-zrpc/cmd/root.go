/*
Copyright © 2024 robin zhang
*/
package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

const (
	version = "v1.0.0"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:     "protoc-gen-zrpc",
	Short:   "A protoc plugin for zrpc",
	Version: version,
	// Run: func(cmd *cobra.Command, args []string) {},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.protoc-gen-zrpc.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	// rootCmd.Flags().BoolVarP(&flagVersion, "version", "V", false, "show version")
}
