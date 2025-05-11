package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"homelab-inventory/pkg/model"
	"net/http"
)

// CheckHealth performs a GET /health request to ensure server is up.
func CheckHealth(endpoint string) error {
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

// PushSystemInfo sends system info to the remote server via POST /sysinfo.
func PushSystemInfo(endpoint string, info *model.SystemInfo) error {
	if err := CheckHealth(endpoint); err != nil {
		return fmt.Errorf("server not healthy: %w", err)
	}

	body, err := json.Marshal(info)
	if err != nil {
		return fmt.Errorf("failed to marshal info: %w", err)
	}

	resp, err := http.Post(endpoint+"/sysinfo", "application/json", bytes.NewBuffer(body))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status: %d", resp.StatusCode)
	}

	return nil
}
