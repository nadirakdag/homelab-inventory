package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"homelab-inventory/internal/collector"
	"homelab-inventory/pkg/model"

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
			if err := postInfo(url, info); err != nil {
				fmt.Println("Failed to send system info:", err)
			} else {
				fmt.Println("System info sent to", url)
			}
		}
	},
}

func postInfo(endpoint string, info *model.SystemInfo) error {
	if err := checkHealth(endpoint); err != nil {
		return fmt.Errorf("server not healthy: %w", err)
	}

	body, _ := json.Marshal(info)
	resp, err := http.Post(endpoint, "application/json", bytes.NewBuffer(body))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status: %d", resp.StatusCode)
	}

	return nil
}

func checkHealth(endpoint string) error {
	resp, err := http.Get(endpoint + "/health")
	if err != nil {
		return fmt.Errorf("health check failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("health check returned %d", resp.StatusCode)
	}
	return nil
}
