package model

type PhysicalDisk struct {
	Name         string  `json:"name"`  // e.g., sda, nvme0n1
	Model        string  `json:"model"` // e.g., Samsung SSD 970 EVO
	SerialNumber string  `json:"serial_number"`
	SizeGB       float64 `json:"size_gb"`   // physical size
	Type         string  `json:"type"`      // HDD, SSD, NVMe, etc.
	Interface    string  `json:"interface"` // SATA, NVMe, USB
	Vendor       string  `json:"vendor"`
}
