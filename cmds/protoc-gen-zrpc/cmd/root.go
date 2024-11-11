/*
Copyright © 2024 robin zhang
*/
package cmd

import (
	"github.com/nicolerobin/zrpc/cmds/protoc-gen-zrpc/constant"
	"github.com/nicolerobin/zrpc/cmds/protoc-gen-zrpc/generator"
	"os"

	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:     "protoc-gen-zrpc",
	Short:   "A protoc plugin for zrpc",
	Version: constant.Version,
	Run: func(cmd *cobra.Command, args []string) {
		generator.Generate()
	},
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
