package scanners

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"

	"swift-anticheat-scanner/pkg"
)

var knownCheatProcesses = []string{
	// --- Minecraft Ghost Clients ---
	"vape", "vapeinjector", "vapeinjection", "vapeclient",
	"vapev4", "vapev3", "vapev2",
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
	"raven", "ravenclient", "ravenb", "ravenbplus",
	"kilo", "kiloclient",
	"dot", "dotclient",
	"loose", "looseclient",

	// --- Free Clients ---
	"forgehax", "forgehaxclient",
	"impact", "impactclient",
	"wurst", "wurstclient",
	"aristois",
	"sigma", "sigmaclient",
	"inertia", "inertiaclient",

	// --- General Cheat Tools ---
	"autoclicker",
	"autoclick",
	"mousekey",
	"opauto",
	"tinytask",
	"pulover",
	"macroscheduler",

	// --- External cheat launchers / injectors ---
	"injection", "injector", "cheatloader", "clientloader",
	"cheatlauncher",
	"dexed", "exe injector", "extremeinjector",
	"processhacker",
	"cheatengine", "cheat engine",
	"artmoney",
	"wemod",

	// --- Ghost clients (paid) ---
	"entropy", "entropyclient",
	"exhibition", "exhibitionclient",
	"mio", "mioclient",
	"cuckold", "cuckoldclient",
	"zero", "zeroclient",
	"fdpclient", "fdp",
	"skilled", "skilledclient",
	"astolfo",
	"drop",
	"azura",
	"blizzard",
	"bubble",
	"celsius",
	"cobalt",
	"corrosion",
	"diablo",
	"dort",
	"dortware",
	"epic",
	"epsilon",
	"ethereal",
	"excuse",
	"flux",
	"fusion",
	"gamesense",
	"gravity",
	"gothaj",
	"huzuni",
	"insanity",
	"interia",
	"jigsaw",
	"kagu",
	"ketamine",
	"konas",
	"koks",
	"lambda",
	"lemon",
	"lithium",
	"lune",
	"lupus",
	"manky",
	"mercury",
	"meteor",
	"monkey",
	"moonlight",
	"nepixel",
	"neverhook",
	"nightly",
	"odyssey",
	"onyx",
	"orbit",
	"ozone",
	"panic",
	"phantom",
	"phobos",
	"pulsive",
	"pyro",
	"reckt",
	"reshack",
	"residual",
	"roofless",
	"rusherhack",
	"ryozan",
	"salhack",
	"sanction",
	"seppuku",
	"shizuku",
	"spicy",
	"spirit",
	"strife",
	"summer",
	"syneid",
	"tifality",
	"toxic",
	"useless",
	"viper",
	"vulcan",
	"winter",
	"xatz",
	"zabex",

	// --- Bypass mods ---
	"anticheatbypass", "acbypass", "bypass",
	"antiblindness", "safetyfirst",

	// --- Utility tools ---
	"xray", "x-ray",
	"wallhack", "esp",
	"killaura", "aura",
	"reach", "reachhack",
	"antibot", "antiknockback",
	"velocityhack", "noslow",
	"bhop", "speedhack",
	"step", "highjump",
	"flyhack", "flight",
	"scaffold", "tower",
	"nuker", "nukah",
	"baritone",
	"tpaura", "aimassist", "aimbot",
	"triggerbot", "autoarmor",
	"autowalk", "invwalker",
	"fucker", "civbreak",

	// --- Config loaders ---
	"prestige", "prestigeclient",
	"moon", "moonclient",
}

func GetSystemInfo() pkg.SystemInfo {
	info := pkg.SystemInfo{
		Hostname:  getHostname(),
		OSVersion: getOSVersion(),
		Username:  getUsername(),
		BootTime:  getBootTime(),
		Uptime:    getUptime(),
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
		results = ScanWindowsProcesses()
	}

	for i, proc := range results {
		results[i].Suspicious = isSuspiciousProcess(proc.Name)
		if results[i].Suspicious {
			results[i].Reason = fmt.Sprintf("Known cheat process: %s", proc.Name)
		}
	}

	return results
}

func ScanServices() []string {
	var suspicious []string

	if runtime.GOOS != "windows" {
		return suspicious
	}

	out, err := exec.Command("powershell", "-NoProfile", "-Command",
		"Get-Service | Where-Object { $_.StartType -ne 'Disabled' } | Select-Object Name, DisplayName, Status | ConvertTo-Json -Compress").Output()
	if err != nil {
		return suspicious
	}

	lines := strings.Split(string(out), "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" || line == "[" || line == "]" {
			continue
		}
		line = strings.TrimPrefix(line, ",")
		lower := strings.ToLower(line)

		for _, cheat := range knownCheatProcesses {
			if strings.Contains(lower, cheat) {
				suspicious = append(suspicious, extractServiceName(line))
				break
			}
		}
	}

	return suspicious
}

func ScanInstalledPrograms() []string {
	var results []string

	if runtime.GOOS != "windows" {
		return results
	}

	out, err := exec.Command("powershell", "-NoProfile", "-Command",
		"Get-ItemProperty 'HKLM:\\Software\\Microsoft\\Windows\\CurrentVersion\\Uninstall\\*', 'HKCU:\\Software\\Microsoft\\Windows\\CurrentVersion\\Uninstall\\*' | Select-Object DisplayName | Where-Object { $_.DisplayName } | ConvertTo-Json -Compress").Output()
	if err != nil {
		return results
	}

	lines := strings.Split(string(out), "\n")
	for _, line := range lines {
		lower := strings.ToLower(line)
		for _, cheat := range knownCheatProcesses {
			if strings.Contains(lower, cheat) {
				results = append(results, extractProgramName(line))
				break
			}
		}
	}

	return results
}

func ScanStartupPrograms() []pkg.StartupEntry {
	var results []pkg.StartupEntry

	if runtime.GOOS != "windows" {
		return results
	}

	out, err := exec.Command("powershell", "-NoProfile", "-Command",
		"Get-CimInstance Win32_StartupCommand | Select-Object Name, Command, Location | ConvertTo-Json -Compress").Output()
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
		lower := strings.ToLower(line)

		suspicious := false
		for _, cheat := range knownCheatProcesses {
			if strings.Contains(lower, cheat) {
				suspicious = true
				break
			}
		}

		name, command := extractStartupInfo(line)
		results = append(results, pkg.StartupEntry{
			Name:       name,
			Command:    command,
			Suspicious: suspicious,
		})
	}

	return results
}

func ScanWindowsProcesses() []pkg.ProcessInfo {
	var results []pkg.ProcessInfo

	out, err := exec.Command("powershell", "-NoProfile", "-Command",
		"Get-Process | Select-Object Id, ProcessName, Path, Company | ConvertTo-Json -Compress").Output()
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

func extractServiceName(line string) string {
	if idx := strings.Index(line, `"Name":`); idx != -1 {
		start := idx + 7
		end := strings.Index(line[start:], "\"")
		if end != -1 {
			return line[start : start+end]
		}
	}
	return line
}

func extractProgramName(line string) string {
	if idx := strings.Index(line, `"DisplayName":`); idx != -1 {
		start := idx + 14
		end := strings.Index(line[start:], "\"")
		if end != -1 {
			return line[start : start+end]
		}
	}
	return line
}

func extractStartupInfo(line string) (string, string) {
	name := ""
	command := ""
	if idx := strings.Index(line, `"Name":`); idx != -1 {
		start := idx + 7
		end := strings.Index(line[start:], "\"")
		if end != -1 {
			name = line[start : start+end]
		}
	}
	if idx := strings.Index(line, `"Command":`); idx != -1 {
		start := idx + 10
		end := strings.Index(line[start:], "\"")
		if end != -1 {
			command = line[start : start+end]
		}
	}
	return name, command
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
