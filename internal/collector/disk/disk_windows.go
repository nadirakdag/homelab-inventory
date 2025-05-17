//go:build windows

package disk

import (
	"homelab-inventory/pkg/model"
	"strings"
)

type Win32DiskDrive struct {
	DeviceID      string
	Model         string
	SerialNumber  string
	Size          uint64
	InterfaceType string
	MediaType     string // Not always populated
}

func CollectPhysicalDisks() ([]model.PhysicalDisk, error) {
	var drives []Win32DiskDrive
	err := wmi.Query("SELECT DeviceID, Model, SerialNumber, Size, InterfaceType, MediaType FROM Win32_DiskDrive", &drives)
	if err != nil {
		return nil, err
	}

	var results []model.PhysicalDisk
	for _, d := range drives {
		disk := model.PhysicalDisk{
			Name:         strings.TrimPrefix(d.DeviceID, `\\.\`),
			Model:        strings.TrimSpace(d.Model),
			SerialNumber: strings.TrimSpace(d.SerialNumber),
			SizeGB:       float64(d.Size) / (1024 * 1024 * 1024), // Bytes to GB
			Type:         resolveDiskType(d.MediaType),
			Interface:    strings.TrimSpace(d.InterfaceType),
		}
		results = append(results, disk)
	}

	return results, nil
}

func resolveDiskType(mediaType string) string {
	mediaType = strings.ToLower(mediaType)
	switch {
	case strings.Contains(mediaType, "ssd"):
		return "SSD"
	case strings.Contains(mediaType, "hdd"):
		return "HDD"
	case mediaType != "":
		return mediaType
	default:
		return "Unknown"
	}
}
