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
	"sigma", "impact", "baritone", "cheatengine", "wemod",
	"artmoney", "processhacker", "extremeinjector", "xenos",
	"dexed", "entropy", "exhibition", "astolfo", "prestige",
	"meteor", "phobos", "rusherhack", "salhack", "seppuku",
	"gamesense", "konas", "koks", "lambda", "lemon",
	"moonlight", "ozone", "pyro", "spicy", "summer",
	"toxic", "vulcan", "winter", "zero",
}

var suspiciousRegistryPaths = []string{
	`HKCU:\Software\Microsoft\Windows\CurrentVersion\Run`,
	`HKCU:\Software\Microsoft\Windows\CurrentVersion\RunOnce`,
	`HKLM:\Software\Microsoft\Windows\CurrentVersion\Run`,
	`HKLM:\Software\Microsoft\Windows\CurrentVersion\RunOnce`,
	`HKCU:\Software\Microsoft\Windows\CurrentVersion\RunServices`,
	`HKCU:\Software\Microsoft\Windows NT\CurrentVersion\Windows`,
	`HKCU:\Software\Microsoft\Windows\CurrentVersion\App Paths`,
}

var extraRegistryChecks = []string{
	`HKCU:\Software\Microsoft\Windows\CurrentVersion\Explorer\RecentDocs`,
	`HKCU:\Software\Microsoft\Windows\CurrentVersion\Explorer\ComDlg32\OpenSavePidlMRU`,
	`HKCU:\Software\Classes\Local Settings\Software\Microsoft\Windows\CurrentVersion\AppContainer\Storage`,
}

func ScanWindowsArtifacts() pkg.WindowsArtifacts {
	artifacts := pkg.WindowsArtifacts{
		PrefetchFiles:          scanPrefetch(),
		SuspiciousRegistryKeys: scanRegistry(),
		SuspiciousRunKeys:      scanRunKeys(),
		PowershellHistory:      getPowershellHistory(),
		EventLogCleared:        checkEventLogCleared(),
		RecentDocuments:        scanRecentDocuments(),
		DnsCache:               scanDnsCache(),
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
			"Get-ItemProperty '"+regPath+"' -ErrorAction SilentlyContinue | Select-Object -ExpandProperty * -ErrorAction SilentlyContinue").Output()
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

	for _, regPath := range extraRegistryChecks {
		out, err := exec.Command("powershell", "-NoProfile", "-Command",
			"Get-ItemProperty '"+regPath+"' -ErrorAction SilentlyContinue | Select-Object -ExpandProperty * -ErrorAction SilentlyContinue").Output()
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

func scanRunKeys() []string {
	var results []string

	runPaths := []string{
		`HKCU:\Software\Microsoft\Windows\CurrentVersion\Run`,
		`HKLM:\Software\Microsoft\Windows\CurrentVersion\Run`,
		`HKCU:\Software\Microsoft\Windows\CurrentVersion\RunOnce`,
		`HKLM:\Software\Microsoft\Windows\CurrentVersion\RunOnce`,
	}

	for _, path := range runPaths {
		out, err := exec.Command("powershell", "-NoProfile", "-Command",
			"Get-ItemProperty '"+path+"' -ErrorAction SilentlyContinue | Select-Object * | Format-List").Output()
		if err != nil {
			continue
		}

		lines := strings.Split(string(out), "\n")
		for _, line := range lines {
			lower := strings.ToLower(line)
			for _, cheat := range knownCheatProcesses {
				if strings.Contains(lower, cheat) {
					results = append(results, strings.TrimSpace(line))
					break
				}
			}
		}
	}

	return results
}

func getPowershellHistory() string {
	historyPaths := []string{
		os.Getenv("USERPROFILE") + "\\AppData\\Roaming\\Microsoft\\Windows\\PowerShell\\PSReadLine\\ConsoleHost_history.txt",
		os.Getenv("USERPROFILE") + "\\AppData\\Roaming\\Microsoft\\Windows\\PowerShell\\PSReadLine\\WindowsPowerShell\\ConsoleHost_history.txt",
		os.Getenv("USERPROFILE") + "\\AppData\\Roaming\\Microsoft\\PowerShell\\PSReadLine\\ConsoleHost_history.txt",
	}

	for _, historyPath := range historyPaths {
		data, err := os.ReadFile(historyPath)
		if err == nil && len(data) > 0 {
			return string(data)
		}
	}

	return ""
}

func checkEventLogCleared() bool {
	out, err := exec.Command("powershell", "-NoProfile", "-Command",
		"Get-WinEvent -FilterHashtable @{LogName='System'; Id=1102,104} -MaxEvents 1 -ErrorAction SilentlyContinue | Select-Object TimeCreated, Id | Format-Table -HideTableHeaders").Output()
	if err != nil {
		return false
	}
	return len(strings.TrimSpace(string(out))) > 0
}

func scanRecentDocuments() []string {
	var results []string

	recentPaths := []string{
		os.Getenv("USERPROFILE") + "\\Recent",
		os.Getenv("APPDATA") + "\\Microsoft\\Windows\\Recent",
		os.Getenv("USERPROFILE") + "\\AppData\\Roaming\\Microsoft\\Windows\\Recent",
	}

	for _, recentPath := range recentPaths {
		if _, err := os.Stat(recentPath); os.IsNotExist(err) {
			continue
		}

		entries, err := os.ReadDir(recentPath)
		if err != nil {
			continue
		}

		for _, entry := range entries {
			if entry.IsDir() {
				continue
			}
			name := strings.ToLower(entry.Name())
			for _, cheat := range suspiciousFileNames {
				if strings.Contains(name, cheat) {
					results = append(results, filepath.Join(recentPath, entry.Name()))
					break
				}
			}
		}
	}

	return results
}

func scanDnsCache() []string {
	var results []string

	out, err := exec.Command("powershell", "-NoProfile", "-Command",
		"Get-DnsClientCache -ErrorAction SilentlyContinue | Where-Object { $_.Entry -match 'cheat|hack|vape|inject|client|ghost|minecraft' } | Select-Object Entry | Format-Table -HideTableHeaders").Output()
	if err != nil {
		return results
	}

	lines := strings.Split(string(out), "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line != "" {
			results = append(results, line)
		}
	}

	return results
}
