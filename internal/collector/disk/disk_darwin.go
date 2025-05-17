//go:build darwin

package disk

import (
	"fmt"
	"homelab-inventory/pkg/model"
	"os/exec"
	"strconv"
	"strings"
)

func CollectPhysicalDisks() ([]model.PhysicalDisk, error) {
	output, err := exec.Command("system_profiler", "SPStorageDataType").Output()
	if err != nil {
		return nil, fmt.Errorf("failed to execute system_profiler: %w", err)
	}

	return parseDiskOutput(string(output)), nil
}

func parseDiskOutput(output string) []model.PhysicalDisk {
	lines := strings.Split(output, "\n")
	var current model.PhysicalDisk
	diskMap := make(map[string]model.PhysicalDisk)

	for _, line := range lines {
		line = strings.TrimSpace(line)

		switch {
		case strings.HasPrefix(line, "Device Name:"):
			current.Name = trimPrefix(line, "Device Name:")

		case strings.HasPrefix(line, "Media Name:"):
			current.Model = trimPrefix(line, "Media Name:")

		case strings.HasPrefix(line, "Medium Type:"):
			current.Type = trimPrefix(line, "Medium Type:")

		case strings.HasPrefix(line, "Capacity:"):
			current.SizeGB = parseSizeGB(line)

		case strings.HasPrefix(line, "Protocol:"):
			current.Interface = trimPrefix(line, "Protocol:")

		case strings.HasPrefix(line, "Vendor:"):
			current.Vendor = trimPrefix(line, "Vendor:")

		case line == "":
			addUniqueDisk(diskMap, current)
			current = model.PhysicalDisk{}
		}
	}

	// Add last parsed block
	if current.Model != "" {
		addUniqueDisk(diskMap, current)
	}

	// Convert map to slice
	var disks []model.PhysicalDisk
	for _, d := range diskMap {
		disks = append(disks, d)
	}
	return disks
}

func addUniqueDisk(diskMap map[string]model.PhysicalDisk, disk model.PhysicalDisk) {
	if disk.Model == "" {
		return
	}

	key := fmt.Sprintf("%s-%s-%.1f-%s", disk.Model, disk.Type, disk.SizeGB, disk.Interface)
	if _, exists := diskMap[key]; !exists {
		diskMap[key] = disk
	}
}

func trimPrefix(line, prefix string) string {
	return strings.TrimSpace(strings.TrimPrefix(line, prefix))
}

func parseSizeGB(line string) float64 {
	parts := strings.Fields(line)
	if len(parts) < 2 {
		return 0
	}

	sizeStr := strings.ReplaceAll(parts[1], ",", "")
	size, err := strconv.ParseFloat(sizeStr, 64)
	if err != nil {
		return 0
	}

	return size
}
