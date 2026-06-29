package main

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"strings"
	"time"

	"swift-anticheat-scanner/pkg"
	"swift-anticheat-scanner/scanners"
)

var Banner = `
  ███████  ██    ██  ██  ███████  ████████  ██
  ██       ██    ██  ██     ██       ██     ██
  ███████  ██    ██  ██    ██        ██     ██
       ██  ██    ██  ██   ██         ██     ██
  ███████   ██████   ██  ███████     ██     ██
  =============================================
  SwiftAntiCheat Scanner v1.0.0
  =============================================
`

func main() {
	scanID := flag.String("id", "", "Scan ID from the server")
	apiURL := flag.String("api", "http://localhost:3000", "API base URL")
	flag.Parse()

	fmt.Println(Banner)
	fmt.Println("SwiftAntiCheat Scanner starting...")
	fmt.Println()

	if *scanID == "" {
		fmt.Println("No scan ID provided. Usage: swiftac-scanner.exe -id YOUR_SCAN_ID")
		fmt.Println("Get your scan ID from the server staff.")
		waitAndExit()
		return
	}

	results := pkg.ScanReport{
		ScanID:       *scanID,
		ScanTime:     time.Now().UTC().Format(time.RFC3339),
		SystemInfo:   pkg.SystemInfo{},
		Processes:    []pkg.ProcessInfo{},
		SuspiciousFiles: []pkg.FileInfo{},
		WindowsArtifacts: pkg.WindowsArtifacts{},
		MinecraftMods: []pkg.MinecraftMod{},
		Flags:        []pkg.Flag{},
	}

	fmt.Println("[1/5] Collecting system information...")
	results.SystemInfo = scanners.GetSystemInfo()
	hwid := generateHWID(results.SystemInfo)
	fmt.Printf("       Hostname: %s\n", results.SystemInfo.Hostname)
	fmt.Printf("       OS: %s\n", results.SystemInfo.OSVersion)

	fmt.Println("[2/5] Scanning running processes...")
	results.Processes = scanners.ScanProcesses()
	for _, proc := range results.Processes {
		if proc.Suspicious {
			results.Flags = append(results.Flags, pkg.Flag{
				Type:     "suspicious_process",
				Severity: "high",
				Name:     proc.Name,
				Detail:   fmt.Sprintf("Suspicious process running: %s (PID: %d)", proc.Name, proc.PID),
			})
		}
	}
	fmt.Printf("       Found %d processes (%d suspicious)\n", len(results.Processes), countSuspicious(results.Processes))

	fmt.Println("[3/5] Scanning filesystem for cheats...")
	results.SuspiciousFiles = scanners.ScanFilesystem()
	for _, f := range results.SuspiciousFiles {
		results.Flags = append(results.Flags, pkg.Flag{
			Type:     "suspicious_file",
			Severity: f.Severity,
			Name:     f.Name,
			Detail:   fmt.Sprintf("Found: %s", f.Path),
		})
	}
	fmt.Printf("       Found %d suspicious files\n", len(results.SuspiciousFiles))

	fmt.Println("[4/5] Checking Windows artifacts...")
	results.WindowsArtifacts = scanners.ScanWindowsArtifacts()
	for _, pf := range results.WindowsArtifacts.PrefetchFiles {
		if pf.Suspicious {
			results.Flags = append(results.Flags, pkg.Flag{
				Type:     "prefetch",
				Severity: "medium",
				Name:     pf.Name,
				Detail:   fmt.Sprintf("Suspicious prefetch entry: %s (last run: %s)", pf.Name, pf.LastRun),
			})
		}
	}
	fmt.Printf("       Checked %d prefetch entries\n", len(results.WindowsArtifacts.PrefetchFiles))

	if len(results.WindowsArtifacts.SuspiciousRegistryKeys) > 0 {
		for _, rk := range results.WindowsArtifacts.SuspiciousRegistryKeys {
			results.Flags = append(results.Flags, pkg.Flag{
				Type:     "registry",
				Severity: "high",
				Name:     rk,
				Detail:   fmt.Sprintf("Suspicious registry key: %s", rk),
			})
		}
	}

	fmt.Println("[5/5] Scanning Minecraft mods...")
	results.MinecraftMods = scanners.ScanMinecraftMods()
	for _, mod := range results.MinecraftMods {
		if mod.Suspicious {
			results.Flags = append(results.Flags, pkg.Flag{
				Type:     "minecraft_mod",
				Severity: "high",
				Name:     mod.Name,
				Detail:   fmt.Sprintf("Suspicious Minecraft mod: %s (%s)", mod.Name, mod.Path),
			})
		}
	}
	fmt.Printf("       Scanned %d mods\n", len(results.MinecraftMods))

	results.HWIDHash = hwid

	fmt.Println()
	fmt.Println("========================================")
	fmt.Println("Scan complete! Uploading results...")
	fmt.Println("========================================")

	err := pkg.UploadResults(*apiURL, results)
	if err != nil {
		fmt.Printf("Failed to upload results: %v\n", err)
		fmt.Println("Results saved locally to scan_result.json")
		saveLocal(results)
		waitAndExit()
		return
	}

	fmt.Println("Results uploaded successfully!")
	fmt.Println("You can now close this window.")
	waitAndExit()
}

func generateHWID(info pkg.SystemInfo) string {
	hwidStr := strings.Join([]string{
		info.Hostname,
		info.MacAddress,
		info.CPUInfo,
		info.MotherboardSerial,
	}, "|")
	hash := sha256.Sum256([]byte(hwidStr))
	return hex.EncodeToString(hash[:])
}

func countSuspicious(procs []pkg.ProcessInfo) int {
	count := 0
	for _, p := range procs {
		if p.Suspicious {
			count++
		}
	}
	return count
}

func saveLocal(results pkg.ScanReport) {
	data, _ := json.MarshalIndent(results, "", "  ")
	os.WriteFile("scan_result.json", data, 0644)
}

func waitAndExit() {
	fmt.Println()
	fmt.Println("Press Enter to exit...")
	fmt.Scanln()
	os.Exit(0)
}
