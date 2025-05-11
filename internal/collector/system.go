package collector

import (
	"fmt"
	"homelab-inventory/pkg/model"
	"runtime"

	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/disk"
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

	partitions, err := disk.Partitions(true)
	if err != nil {
		return nil, fmt.Errorf("listing partitions: %w", err)
	}

	var disks []model.DiskInfo
	for _, p := range partitions {
		usage, err := disk.Usage(p.Mountpoint)
		if err != nil || usage.Total == 0 {
			continue
		}

		disks = append(disks, model.DiskInfo{
			Mountpoint: p.Mountpoint,
			TotalGB:    float64(usage.Total) / (1024 * 1024 * 1024),
			UsedGB:     float64(usage.Used) / (1024 * 1024 * 1024),
			FreeGB:     float64(usage.Free) / (1024 * 1024 * 1024),
		})
	}

	info := &model.SystemInfo{
		Hostname: hostInfo.Hostname,
		OS:       hostInfo.OS,
		Platform: hostInfo.Platform,
		Arch:     runtime.GOARCH,
		CPUModel: cpuInfo[0].ModelName,
		CPUCores: runtime.NumCPU(),
		MemoryGB: float64(vmStat.Total) / (1024 * 1024 * 1024),
		Disks:    disks,
	}

	return info, nil
}
