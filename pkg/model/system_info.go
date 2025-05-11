package model

type DiskInfo struct {
	Mountpoint string  `json:"mountpoint"`
	TotalGB    float64 `json:"total_gb"`
	UsedGB     float64 `json:"used_gb"`
	FreeGB     float64 `json:"free_gb"`
}

type SystemInfo struct {
	Hostname string     `json:"hostname"`
	OS       string     `json:"os"`
	Platform string     `json:"platform"`
	Arch     string     `json:"arch"`
	CPUModel string     `json:"cpu_model"`
	CPUCores int        `json:"cpu_cores"`
	MemoryGB float64    `json:"memory_gb"`
	Disks    []DiskInfo `json:"disks"`
}
