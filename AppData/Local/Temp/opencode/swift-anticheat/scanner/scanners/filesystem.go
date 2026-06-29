package scanners

import (
	"os"
	"path/filepath"
	"strings"
	"time"

	"swift-anticheat-scanner/pkg"
)

var suspiciousFileNames = []string{
	"vape", "liquidbounce", "future", "rise", "novoline",
	"whiteout", "cracked", "raven", "aero", "loose",
	"injection", "injector", "cheatloader", "clientloader",
	"anticheatbypass", "killaura", "autoclicker", "reach",
	"xray", "x-ray", "esp", "wallhack", "nodus",
	"wurst", "aristois", "sigma", "impact", "baritone",
}

var suspiciousExtensions = []string{
	".exe", ".dll", ".jar", ".bat", ".ps1", ".vbs", ".scr",
}

var scanPaths = []string{
	os.Getenv("USERPROFILE") + "\\Desktop",
	os.Getenv("USERPROFILE") + "\\Downloads",
	os.Getenv("TEMP"),
	os.Getenv("LOCALAPPDATA") + "\\Temp",
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
				reason = "Executable jar in download area"
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
		"whiteout", "ravenb+", "cheat", "hack", "client",
	}
	for _, d := range cheatDirs {
		if strings.Contains(lower, d) {
			return true
		}
	}
	return false
}
