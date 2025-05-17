package model

type SystemInfo struct {
	Hostname      string         `json:"hostname"`
	OS            string         `json:"os"`
	Platform      string         `json:"platform"`
	Arch          string         `json:"arch"`
	CPUModel      string         `json:"cpu_model"`
	CPUCores      int            `json:"cpu_cores"`
	MemoryGB      float64        `json:"memory_gb"`
	PhysicalDisks []PhysicalDisk `json:"physical_disks"`
}
