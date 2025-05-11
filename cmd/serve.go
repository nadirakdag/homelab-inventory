package cmd

import (
	"homelab-inventory/internal/server"

	"github.com/spf13/cobra"
)

var port string

func init() {
	serveCmd.Flags().StringVar(&port, "port", "8080", "API server port")
	rootCmd.AddCommand(serveCmd)
}

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Start system info API server",
	Run: func(cmd *cobra.Command, args []string) {
		server.StartServer(port)
	},
}
