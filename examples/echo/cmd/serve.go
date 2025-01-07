/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"context"

	"github.com/nicolerobin/zrpc/log"
	"github.com/nicolerobin/zrpc/server"
	"github.com/spf13/cobra"
	"go.uber.org/zap"

	"echo/handler"
	pb "echo/api/echo"
)

func ServeCmd(cmd *cobra.Command, args []string) {
	ctx := context.Background()
	server.RegisterService(pb.GreeterService, handler.GreeterHandler{})
}

// serveCmd represents the serve command
var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Run server",
	Run:   ServeCmd,
}

func init() {
	rootCmd.AddCommand(serveCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// serveCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// serveCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
