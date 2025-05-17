//go:build linux

package disk

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"

	"homelab-inventory/internal/logging"
	"homelab-inventory/pkg/model"
)

func CollectPhysicalDisks() ([]model.PhysicalDisk, error) {
	output, err := exec.Command("lsblk", "-d", "-o", "NAME,SIZE,MODEL,TYPE", "-n").Output()
	if err != nil {
		return nil, fmt.Errorf("failed to execute lsblk: %w", err)
	}

	return parseLSBLKOutput(string(output)), nil
}

func parseLSBLKOutput(output string) []model.PhysicalDisk {
	var disks []model.PhysicalDisk
	lines := strings.Split(output, "\n")
	logging.Logger.Debugf("output", "output", output)
	logging.Logger.Debugf("Parsing lsblk output", "lines", lines)

	for _, line := range lines {
		line = strings.TrimSpace(line)
		logging.Logger.Debugf("Line", "line", line)
		if line == "" {
			continue
		}

		fields := strings.Fields(line)
		logging.Logger.Debugf("Fields", "fields", fields)
		if fields[len(fields)-1] != "disk" {
			logging.Logger.Debugf("Skipping line", "line", line)
			continue
		}

		name := fields[0]
		logging.Logger.Debugf("Name", "name", name)

		sizeGB := parseSizeGB(fields[1])
		logging.Logger.Debugf("Size", "sizeGB", sizeGB)

		modelType := strings.Join(fields[2:len(fields)-1], " ")
		logging.Logger.Debugf("Model", "modelType", modelType)
		if modelType == "" {
			modelType = guessModelFromDevice(name) // Fallback for ARM devices
		}

		disk := model.PhysicalDisk{
			Name:         name,
			Model:        modelType,
			SerialNumber: fallbackString(readSysFS(fmt.Sprintf("/sys/block/%s/device/serial", name)), "Unavailable"),
			Vendor:       fallbackString(readSysFS(fmt.Sprintf("/sys/block/%s/device/vendor", name)), "Unknown"),
			Interface:    detectInterface(name),
			Type:         resolveDiskType(name),
			SizeGB:       sizeGB,
		}

		disks = append(disks, disk)
	}

	return disks
}

func fallbackString(value, fallback string) string {
	if strings.TrimSpace(value) == "" {
		return fallback
	}
	return value
}

func guessModelFromDevice(name string) string {
	switch {
	case strings.HasPrefix(name, "mmcblk"):
		return "Embedded eMMC"
	case strings.HasPrefix(name, "sd"):
		return "SD Card"
	case strings.HasPrefix(name, "nvme"):
		return "NVMe Device"
	case strings.HasPrefix(name, "usb"):
		return "USB Storage"
	default:
		return "Unknown"
	}
}

func parseSizeGB(raw string) float64 {
	value, err := strconv.ParseFloat(strings.TrimSuffix(raw, "G"), 64)
	if err != nil {
		return 0
	}
	return value
}

func resolveDiskType(dev string) string {
	if isSSD(dev) {
		return "SSD"
	}
	return "HDD"
}

func isSSD(dev string) bool {
	rotPath := fmt.Sprintf("/sys/block/%s/queue/rotational", dev)
	data, err := os.ReadFile(rotPath)
	return err == nil && strings.TrimSpace(string(data)) == "0"
}

func detectInterface(dev string) string {
	link, err := filepath.EvalSymlinks(fmt.Sprintf("/sys/block/%s", dev))
	if err != nil {
		return "unknown"
	}

	for _, part := range strings.Split(link, "/") {
		switch {
		case strings.HasPrefix(part, "usb"):
			return "USB"
		case strings.HasPrefix(part, "ata"):
			return "SATA"
		case strings.HasPrefix(part, "nvme"):
			return "NVMe"
		case strings.HasPrefix(part, "mmc"):
			return "eMMC"
		}
	}

	return "unknown"
}

func readSysFS(path string) string {
	data, err := os.ReadFile(path)
	if err != nil {
		return ""
	}
	return strings.TrimSpace(string(data))
}
