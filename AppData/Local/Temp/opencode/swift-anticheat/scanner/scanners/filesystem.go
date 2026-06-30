package scanners

import (
	"os"
	"path/filepath"
	"strings"
	"time"

	"swift-anticheat-scanner/pkg"
)

var suspiciousFileNames = []string{
	// Minecraft ghost clients
	"vape", "liquidbounce", "future", "rise", "novoline",
	"whiteout", "cracked", "raven", "aero", "loose",
	"entropy", "exhibition", "mio", "cuckold", "zero",
	"fdpclient", "skilled", "astolfo", "drop", "azura",
	"blizzard", "bubble", "celsius", "cobalt", "corrosion",
	"diablo", "dort", "epic", "epsilon", "ethereal",
	"excuse", "flux", "fusion", "gamesense", "gravity",
	"gothaj", "huzuni", "insanity", "interia", "jigsaw",
	"kagu", "ketamine", "konas", "koks", "lambda", "lemon",
	"lithium", "lune", "lupus", "manky", "mercury",
	"meteor", "monkey", "moonlight", "nepixel", "neverhook",
	"nightly", "odyssey", "onyx", "orbit", "ozone", "panic",
	"phantom", "phobos", "prestige", "pulsive", "pyro",
	"reckt", "reshack", "residual", "roofless", "rusherhack",
	"ryozan", "salhack", "sanction", "seppuku", "shizuku",
	"spicy", "spirit", "strife", "summer", "syneid",
	"tifality", "toxic", "useless", "viper", "vulcan",
	"winter", "xatz", "zabex",

	// Cheat-related terms
	"injection", "injector", "cheatloader", "clientloader",
	"anticheatbypass", "acbypass", "bypass",
	"killaura", "autoclicker", "reach", "triggerbot",
	"xray", "x-ray", "esp", "wallhack", "aimbot", "aimassist",
	"nodus", "wurst", "aristois", "sigma", "impact", "baritone",
	"bhop", "speedhack", "flyhack", "scaffold", "tower",
	"nuker", "antibot", "velocity", "antikb",

	// General hacking tools
	"cheatengine", "cheat engine", "cheatengine.exe",
	"cheatengine-x86_64", "ceserver",
	"artmoney", "artmoney.exe",
	"wemod", "wemod.exe",
	"processhacker", "process hacker",
	"extremeinjector", "extreme injector",
	"dexed", "dexed.exe",
	"xenos", "xenos.exe",
	"injector", "winject", "winject.exe",
}

var suspiciousExtensions = []string{
	".exe", ".dll", ".jar", ".bat", ".ps1", ".vbs", ".scr",
	".cmd", ".js", ".jse", ".wsf", ".wsh", ".msi", ".vxd",
	".sys", ".com", ".pif", ".gadget", ".application",
}

var scanPaths = []string{
	// User profile folders
	os.Getenv("USERPROFILE") + "\\Desktop",
	os.Getenv("USERPROFILE") + "\\Downloads",
	os.Getenv("USERPROFILE") + "\\Documents",
	os.Getenv("USERPROFILE") + "\\OneDrive\\Desktop",
	os.Getenv("USERPROFILE") + "\\OneDrive\\Downloads",
	os.Getenv("USERPROFILE") + "\\OneDrive\\Documents",

	// Temp directories
	os.Getenv("TEMP"),
	os.Getenv("LOCALAPPDATA") + "\\Temp",
	os.Getenv("SYSTEMROOT") + "\\Temp",
	os.Getenv("WINDIR") + "\\Temp",

	// AppData cheat hiding spots
	os.Getenv("APPDATA"),
	os.Getenv("LOCALAPPDATA"),
	os.Getenv("USERPROFILE") + "\\AppData\\Roaming",
	os.Getenv("USERPROFILE") + "\\AppData\\Local",

	// Common game-related data directories
	os.Getenv("APPDATA") + "\\.minecraft",
	os.Getenv("APPDATA") + "\\..\\Roaming\\.minecraft",
	os.Getenv("LOCALAPPDATA") + "\\Programs",
	os.Getenv("LOCALAPPDATA") + "\\Microsoft\\Windows\\INetCache",

	// Startup folder
	os.Getenv("APPDATA") + "\\Microsoft\\Windows\\Start Menu\\Programs\\Startup",
	os.Getenv("PROGRAMDATA") + "\\Microsoft\\Windows\\Start Menu\\Programs\\Startup",

	// Program Files (for injected DLLs, etc.)
	os.Getenv("PROGRAMFILES") + "\\Java",
	os.Getenv("PROGRAMFILES(X86)") + "\\Java",
}

func ScanFilesystem() []pkg.FileInfo {
	var results []pkg.FileInfo
	seen := make(map[string]bool)

	for _, root := range scanPaths {
		if root == "" {
			continue
		}

		filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
			if err != nil || info == nil || info.IsDir() {
				return nil
			}

			if seen[path] {
				return nil
			}
			seen[path] = true

			lower := strings.ToLower(info.Name())
			ext := strings.ToLower(filepath.Ext(info.Name()))

			if !isSuspiciousExtension(ext) && !isSuspiciousName(lower) {
				return nil
			}

			severity := "medium"
			reason := ""

			if isSuspiciousName(lower) {
				severity = "high"
				reason = "Known cheat identifier in filename"
			} else if ext == ".jar" {
				severity = "medium"
				reason = "Executable jar in scan area"
			} else if ext == ".exe" || ext == ".dll" {
				severity = "medium"
				reason = "Suspicious executable binary"
			}

			if pathContainsCheatFolder(path) {
				severity = "high"
				reason = "Located in cheat-associated directory"
			}

			results = append(results, pkg.FileInfo{
				Name:     info.Name(),
				Path:     path,
				Size:     info.Size(),
				Modified: info.ModTime().Format(time.RFC3339),
				Severity: severity,
				Reason:   reason,
			})

			return nil
		})
	}

	return results
}

func isSuspiciousExtension(ext string) bool {
	for _, e := range suspiciousExtensions {
		if ext == e {
			return true
		}
	}
	return false
}

func isSuspiciousName(name string) bool {
	for _, s := range suspiciousFileNames {
		if strings.Contains(name, s) {
			return true
		}
	}
	return false
}

func pathContainsCheatFolder(path string) bool {
	lower := strings.ToLower(path)
	cheatDirs := []string{
		"vape", "liquidbounce", "future", "rise", "novoline",
		"whiteout", "ravenb+", "cheat", "hack",
		"cheatengine", "wemod", "artmoney",
		"injector", "ghostclient",
		"ghost client", "prestige", "prestigeclient",
	}
	for _, d := range cheatDirs {
		if strings.Contains(lower, d) {
			return true
		}
	}
	return false
}
