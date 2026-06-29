package pkg

type ScanReport struct {
	ScanID           string            `json:"scan_id"`
	ScanTime         string            `json:"scan_time"`
	SystemInfo       SystemInfo        `json:"system_info"`
	Processes        []ProcessInfo     `json:"processes"`
	SuspiciousFiles  []FileInfo        `json:"suspicious_files"`
	WindowsArtifacts WindowsArtifacts  `json:"windows_artifacts"`
	MinecraftMods    []MinecraftMod    `json:"minecraft_mods"`
	Flags            []Flag            `json:"flags"`
	HWIDHash         string            `json:"hwid_hash"`
}

type SystemInfo struct {
	Hostname          string `json:"hostname"`
	OSVersion         string `json:"os_version"`
	CPUInfo           string `json:"cpu_info"`
	TotalRAM          string `json:"total_ram"`
	MacAddress        string `json:"mac_address"`
	MotherboardSerial string `json:"motherboard_serial"`
	Username          string `json:"username"`
	BootTime          string `json:"boot_time"`
	Uptime            string `json:"uptime"`
}

type ProcessInfo struct {
	PID         int    `json:"pid"`
	Name        string `json:"name"`
	Path        string `json:"path,omitempty"`
	Suspicious  bool   `json:"suspicious"`
	Reason      string `json:"reason,omitempty"`
}

type FileInfo struct {
	Name     string `json:"name"`
	Path     string `json:"path"`
	Size     int64  `json:"size"`
	Modified string `json:"modified"`
	Severity string `json:"severity"`
	Reason   string `json:"reason"`
}

type WindowsArtifacts struct {
	PrefetchFiles         []PrefetchEntry `json:"prefetch_files"`
	SuspiciousRegistryKeys []string       `json:"suspicious_registry_keys"`
	RecycleBinRecent      []string        `json:"recycle_bin_recent"`
	PowershellHistory     string          `json:"powershell_history"`
	EventLogCleared       bool            `json:"event_log_cleared"`
}

type PrefetchEntry struct {
	Name       string `json:"name"`
	Path       string `json:"path"`
	LastRun    string `json:"last_run"`
	RunCount   int    `json:"run_count"`
	Suspicious bool   `json:"suspicious"`
}

type MinecraftMod struct {
	Name        string `json:"name"`
	Path        string `json:"path"`
	Size        int64  `json:"size"`
	Suspicious  bool   `json:"suspicious"`
	Reason      string `json:"reason,omitempty"`
}

type Flag struct {
	Type     string `json:"type"`
	Severity string `json:"severity"`
	Name     string `json:"name"`
	Detail   string `json:"detail"`
}
