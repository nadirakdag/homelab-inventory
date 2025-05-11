package cmd

import (
	"fmt"
	"homelab-inventory/internal/version"

	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print version information",
	Run: func(cmd *cobra.Command, args []string) {
		v := version.Get()
		fmt.Printf("Version:    %s\n", v.Version)
		fmt.Printf("Commit:     %s\n", v.Commit)
		fmt.Printf("Build Time: %s\n", v.BuildTime)
		fmt.Printf("Go Version: %s\n", v.GoVersion)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
