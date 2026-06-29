package scanners

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"swift-anticheat-scanner/pkg"
)

var suspiciousPrefetchNames = []string{
	"vape", "liquidbounce", "future", "rise", "novoline",
	"whiteout", "raven", "cheat", "injector", "injection",
	"autoclicker", "killaura", "xray", "wurst", "aristois",
}

var suspiciousRegistryPaths = []string{
	`HKCU:\Software\Microsoft\Windows\CurrentVersion\Run\`,
	`HKCU:\Software\Microsoft\Windows\CurrentVersion\RunOnce\`,
}

func ScanWindowsArtifacts() pkg.WindowsArtifacts {
	artifacts := pkg.WindowsArtifacts{
		PrefetchFiles:         scanPrefetch(),
		SuspiciousRegistryKeys: scanRegistry(),
		PowershellHistory:     getPowershellHistory(),
		EventLogCleared:       checkEventLogCleared(),
	}

	return artifacts
}

func scanPrefetch() []pkg.PrefetchEntry {
	var results []pkg.PrefetchEntry
	prefetchDir := os.Getenv("WINDIR") + "\\Prefetch"

	if _, err := os.Stat(prefetchDir); os.IsNotExist(err) {
		return results
	}

	entries, err := os.ReadDir(prefetchDir)
	if err != nil {
		return results
	}

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		name := strings.ToLower(entry.Name())
		if !strings.HasSuffix(name, ".pf") {
			continue
		}

		baseName := strings.TrimSuffix(name, ".pf")
		baseName = strings.TrimSuffix(baseName, ".exe")
		suspicious := false

		for _, s := range suspiciousPrefetchNames {
			if strings.Contains(baseName, s) {
				suspicious = true
				break
			}
		}

		info, err := entry.Info()
		if err != nil {
			continue
		}

		results = append(results, pkg.PrefetchEntry{
			Name:       baseName,
			Path:       filepath.Join(prefetchDir, entry.Name()),
			LastRun:    info.ModTime().Format("2006-01-02 15:04:05"),
			RunCount:   0,
			Suspicious: suspicious,
		})
	}

	return results
}

func scanRegistry() []string {
	var results []string

	for _, regPath := range suspiciousRegistryPaths {
		out, err := exec.Command("powershell", "-NoProfile", "-Command",
			"Get-ItemProperty '"+regPath+"' | Select-Object -ExpandProperty *").Output()
		if err != nil {
			continue
		}

		lines := strings.Split(string(out), "\n")
		for _, line := range lines {
			lower := strings.ToLower(line)
			for _, s := range suspiciousFileNames {
				if strings.Contains(lower, s) {
					results = append(results, regPath+" -> "+strings.TrimSpace(line))
					break
				}
			}
		}
	}

	return results
}

func getPowershellHistory() string {
	historyPath := os.Getenv("USERPROFILE") + "\\AppData\\Roaming\\Microsoft\\Windows\\PowerShell\\PSReadLine\\ConsoleHost_history.txt"
	data, err := os.ReadFile(historyPath)
	if err != nil {
		return ""
	}
	return string(data)
}

func checkEventLogCleared() bool {
	out, err := exec.Command("powershell", "-NoProfile", "-Command",
		"Get-WinEvent -FilterHashtable @{LogName='System'; Id=1102,104} -MaxEvents 1 | Select-Object TimeCreated | Format-Table -HideTableHeaders").Output()
	if err != nil {
		return false
	}
	return len(strings.TrimSpace(string(out))) > 0
}
