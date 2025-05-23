package collector

import (
	"fmt"
	"homelab-inventory/internal/collector/disk"
	"homelab-inventory/pkg/model"
	"runtime"

	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/host"
	"github.com/shirou/gopsutil/v3/mem"
)

func CollectSystemInfo() (*model.SystemInfo, error) {
	hostInfo, err := host.Info()
	if err != nil {
		return nil, fmt.Errorf("getting host info: %w", err)
	}

	cpuInfo, err := cpu.Info()
	if err != nil {
		return nil, fmt.Errorf("getting cpu info: %w", err)
	}

	vmStat, err := mem.VirtualMemory()
	if err != nil {
		return nil, fmt.Errorf("getting memory info: %w", err)
	}

	physicalDisks, err := disk.CollectPhysicalDisks()
	if err != nil {
		return nil, fmt.Errorf("listing partitions: %w", err)
	}

	info := &model.SystemInfo{
		Hostname:      hostInfo.Hostname,
		OS:            hostInfo.OS,
		Platform:      hostInfo.Platform,
		Arch:          runtime.GOARCH,
		CPUModel:      cpuInfo[0].ModelName,
		CPUCores:      runtime.NumCPU(),
		MemoryGB:      float64(vmStat.Total) / (1024 * 1024 * 1024),
		PhysicalDisks: physicalDisks,
	}

	return info, nil
}
