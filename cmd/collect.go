package cmd

import (
	"encoding/json"
	"homelab-inventory/internal/client"
	"homelab-inventory/internal/collector"
	"homelab-inventory/internal/logging"

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
			logging.Logger.Errorw("Error collecting system info", "error", err)
			return
		}

		data, _ := json.MarshalIndent(info, "", "  ")
		logging.Logger.Infow("Collected system info", "info", string(data))

		if send && url != "" {
			if err := client.PushSystemInfo(url, info); err != nil {
				logging.Logger.Errorw("Error sending system info", "error", err)
			} else {
				logging.Logger.Infow("Sent system info", "info", string(data))
			}
		}
	},
}
