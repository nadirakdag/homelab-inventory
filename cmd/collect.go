package cmd

import (
	"encoding/json"
	"fmt"
	"homelab-inventory/internal/client"
	"homelab-inventory/internal/collector"

	"github.com/spf13/cobra"
)

var send bool
var url string

func init() {
	collectCmd.Flags().BoolVar(&send, "send", false, "Send system info to API server")
	collectCmd.Flags().StringVar(&url, "url", "", "API server URL to send data")
	rootCmd.AddCommand(collectCmd)
}

var collectCmd = &cobra.Command{
	Use:   "collect",
	Short: "Collect and optionally send system information",
	Run: func(cmd *cobra.Command, args []string) {
		info, err := collector.CollectSystemInfo()
		if err != nil {
			fmt.Println("Error collecting system info:", err)
			return
		}

		data, _ := json.MarshalIndent(info, "", "  ")
		fmt.Println("Collected System Info:\n", string(data))

		if send && url != "" {
			if err := client.PushSystemInfo(url, info); err != nil {
				fmt.Println("Failed to send system info:", err)
			} else {
				fmt.Println("System info sent to", url)
			}
		}
	},
}
