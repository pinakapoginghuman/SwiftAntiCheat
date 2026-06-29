package scanners

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"time"

	"swift-anticheat-scanner/pkg"
)

var knownCheatProcesses = []string{
	"vape", "vapeinjector", "vapeinjection", "vapeclient",
	"liquidbounce", "liquidbounceinjector",
	"futureclient", "future",
	"riseclient", "rise",
	"novoline", "novo",
	"pulse", "pulseclient",
	"whiteout", "whiteoutclient",
	"cracked", "crackedclient",
	"aero", "aeroclient",
	"chi", "chiclient",
	"dream", "dreamclient",
	"void", "voidclient",
	"tenacity", "tenacityclient",
	"raven", "ravenclient",
	"ravenb", "ravenbplus",
	"kilo", "kiloclient",
	"dot", "dotclient",
	"loose", "looseclient",
	"cheatbreaker", "cheatbreakerloader",
	"badlion", "badlionclient",
	"lunarclient",
}

func GetSystemInfo() pkg.SystemInfo {
	info := pkg.SystemInfo{
		Hostname: getHostname(),
		OSVersion: getOSVersion(),
		Username: getUsername(),
		BootTime: getBootTime(),
		Uptime:   getUptime(),
	}

	if runtime.GOOS == "windows" {
		info.CPUInfo = getWMICData("cpu", "Name")
		info.MacAddress = getWMICData("nic", "MACAddress")
		info.MotherboardSerial = getWMICData("baseboard", "SerialNumber")
		info.TotalRAM = getWMICData("memphysical", "TotalPhysicalMemory")
	}

	return info
}

func ScanProcesses() []pkg.ProcessInfo {
	var results []pkg.ProcessInfo

	if runtime.GOOS == "windows" {
		results = scanWindowsProcesses()
	}

	return results
}

func scanWindowsProcesses() []pkg.ProcessInfo {
	var results []pkg.ProcessInfo

	out, err := exec.Command("powershell", "-NoProfile", "-Command",
		"Get-Process | Select-Object Id, ProcessName, Path | ConvertTo-Json -Compress").Output()
	if err != nil {
		return results
	}

	lines := strings.Split(string(out), "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" || line == "[" || line == "]" {
			continue
		}
		line = strings.TrimPrefix(line, ",")

		proc := parseProcessJSON(line)
		if proc != nil {
			proc.Suspicious = isSuspiciousProcess(proc.Name)
			results = append(results, *proc)
		}
	}

	return results
}

func parseProcessJSON(line string) *pkg.ProcessInfo {
	line = strings.TrimSpace(line)
	if !strings.HasPrefix(line, "{") {
		return nil
	}

	var proc pkg.ProcessInfo

	if idx := strings.Index(line, `"Id":`); idx != -1 {
		end := strings.Index(line[idx:], ",")
		if end == -1 {
			end = strings.Index(line[idx:], "}")
		}
		if end != -1 {
			idStr := strings.TrimSpace(line[idx+5 : idx+end])
			idStr = strings.TrimSuffix(idStr, ",")
			idStr = strings.TrimSuffix(idStr, "}")
			idStr = strings.TrimSpace(idStr)
			if idStr != "" {
				var pid int
				if _, err := fmt.Sscanf(idStr, "%d", &pid); err == nil {
					proc.PID = pid
				}
			}
		}
	}

	if idx := strings.Index(line, `"ProcessName":`); idx != -1 {
		start := idx + 14
		end := strings.Index(line[start:], "\"")
		if end != -1 {
			proc.Name = strings.ToLower(line[start : start+end])
		}
	}

	if idx := strings.Index(line, `"Path":`); idx != -1 {
		start := idx + 7
		end := strings.Index(line[start:], "\"")
		if end != -1 {
			proc.Path = line[start : start+end]
		}
	}

	return &proc
}

func isSuspiciousProcess(name string) bool {
	name = strings.ToLower(strings.TrimSpace(name))
	if name == "" {
		return false
	}
	for _, cheat := range knownCheatProcesses {
		if strings.Contains(name, cheat) {
			return true
		}
	}
	return false
}

func getHostname() string {
	name, _ := os.Hostname()
	return name
}

func getOSVersion() string {
	if runtime.GOOS == "windows" {
		out, err := exec.Command("powershell", "-NoProfile", "-Command",
			"(Get-WmiObject Win32_OperatingSystem).Caption").Output()
		if err == nil {
			return strings.TrimSpace(string(out))
		}
	}
	return runtime.GOOS
}

func getUsername() string {
	return os.Getenv("USERNAME")
}

func getWMICData(category, property string) string {
	out, err := exec.Command("powershell", "-NoProfile", "-Command",
		"(Get-WmiObject Win32_"+category+")."+property).Output()
	if err != nil {
		return ""
	}
	return strings.TrimSpace(string(out))
}

func getBootTime() string {
	out, err := exec.Command("powershell", "-NoProfile", "-Command",
		"(gcim Win32_OperatingSystem).LastBootUpTime").Output()
	if err != nil {
		return ""
	}
	return strings.TrimSpace(string(out))
}

func getUptime() string {
	out, err := exec.Command("powershell", "-NoProfile", "-Command",
		"(Get-Date) - (gcim Win32_OperatingSystem).LastBootUpTime | Select-Object Days,Hours,Minutes | ConvertTo-Json").Output()
	if err != nil {
		return ""
	}
	return strings.TrimSpace(string(out))
}


