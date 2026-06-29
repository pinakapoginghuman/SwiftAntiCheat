package scanners

import (
	"crypto/md5"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"swift-anticheat-scanner/pkg"
)

var suspiciousModNames = []string{
	"vape", "liquidbounce", "future", "rise", "novoline",
	"whiteout", "raven", "aero", "loose", "kilo",
	"dream", "void", "tenacity", "chi", "pulse",
	"cracked", "injection", "injector", "cheat", "hack",
	"autoclicker", "killaura", "xray", "x-ray", "esp",
	"wallhack", "reach", "triggerbot", "aimbot", "aimassist",
	"antibot", "esp2d", "tracers", "nametags",
}

var modFolderPaths = []string{
	os.Getenv("APPDATA") + "\\.minecraft\\mods",
	os.Getenv("APPDATA") + "\\.minecraft\\versions",
	os.Getenv("APPDATA") + "\\.minecraft\\libraries",
	os.Getenv("APPDATA") + "\\.minecraft\\saves",
}

func ScanMinecraftMods() []pkg.MinecraftMod {
	var results []pkg.MinecraftMod

	for _, mcPath := range modFolderPaths {
		if _, err := os.Stat(mcPath); os.IsNotExist(err) {
			continue
		}

		filepath.Walk(mcPath, func(path string, info os.FileInfo, err error) error {
			if err != nil || info == nil || info.IsDir() {
				return nil
			}

			ext := strings.ToLower(filepath.Ext(info.Name()))
			if ext != ".jar" && ext != ".zip" && ext != ".litemod" {
				return nil
			}

			lower := strings.ToLower(info.Name())
			suspicious := false
			reason := ""

			for _, s := range suspiciousModNames {
				if strings.Contains(lower, s) {
					suspicious = true
					reason = fmt.Sprintf("Known cheat mod: %s", s)
					break
				}
			}

			if !suspicious && isSuspiciousJarSize(info.Size()) {
				suspicious = true
				reason = "Suspicious jar file size"
			}

			results = append(results, pkg.MinecraftMod{
				Name:       info.Name(),
				Path:       path,
				Size:       info.Size(),
				Suspicious: suspicious,
				Reason:     reason,
			})

			return nil
		})
	}

	return results
}

func isSuspiciousJarSize(size int64) bool {
	return size > 0 && size < 102400
}

func hashFile(path string) string {
	data, err := os.ReadFile(path)
	if err != nil {
		return ""
	}
	return fmt.Sprintf("%x", md5.Sum(data))
}
