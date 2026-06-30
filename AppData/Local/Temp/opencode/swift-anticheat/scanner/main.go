package main

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math/big"
	"os"
	"strings"
	"time"

	"swift-anticheat-scanner/pkg"
	"swift-anticheat-scanner/scanners"
)

var Banner = `
  ███████  ██       ██  ██  ███████  ████████
  ██       ██       ██  ██  ██          ██
  ███████  ██   █   ██  ██  █████       ██
       ██  ██  ██   ██  ██  ██          ██
  ███████   █████ ███   ██  ██          ██
  =============================================
  SwiftAntiCheat Scanner v1.1.0
  =============================================
`

func main() {
	fmt.Println(Banner)
	fmt.Println("SwiftAntiCheat Scanner starting...")
	fmt.Println()

	results := pkg.ScanReport{
		ScanTime:           time.Now().UTC().Format(time.RFC3339),
		SystemInfo:         pkg.SystemInfo{},
		Processes:          []pkg.ProcessInfo{},
		SuspiciousFiles:    []pkg.FileInfo{},
		WindowsArtifacts:   pkg.WindowsArtifacts{},
		MinecraftMods:      []pkg.MinecraftMod{},
		StartupPrograms:    []pkg.StartupEntry{},
		InstalledPrograms:  []string{},
		SuspiciousServices: []string{},
		Flags:              []pkg.Flag{},
	}

	fmt.Println("[1/7] Collecting system information...")
	results.SystemInfo = scanners.GetSystemInfo()
	hwid := generateHWID(results.SystemInfo)
	fmt.Printf("       Hostname: %s\n", results.SystemInfo.Hostname)
	fmt.Printf("       OS: %s\n", results.SystemInfo.OSVersion)

	fmt.Println("[2/7] Scanning running processes...")
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

	fmt.Println("[3/7] Scanning filesystem for cheats...")
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

	fmt.Println("[4/7] Checking Windows artifacts...")
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

	if len(results.WindowsArtifacts.SuspiciousRunKeys) > 0 {
		for _, rk := range results.WindowsArtifacts.SuspiciousRunKeys {
			results.Flags = append(results.Flags, pkg.Flag{
				Type:     "startup",
				Severity: "high",
				Name:     rk,
				Detail:   fmt.Sprintf("Suspicious startup entry: %s", rk),
			})
		}
	}

	if results.WindowsArtifacts.EventLogCleared {
		results.Flags = append(results.Flags, pkg.Flag{
			Type:     "event_log_cleared",
			Severity: "high",
			Name:     "Event Log Cleared",
			Detail:   "Windows event log was cleared — common anti-forensic technique",
		})
	}

	if len(results.WindowsArtifacts.DnsCache) > 0 {
		for _, dns := range results.WindowsArtifacts.DnsCache {
			results.Flags = append(results.Flags, pkg.Flag{
				Type:     "dns_cache",
				Severity: "medium",
				Name:     dns,
				Detail:   fmt.Sprintf("Cheat-related DNS cache entry: %s", dns),
			})
		}
	}

	fmt.Println("[5/7] Scanning startup programs...")
	results.StartupPrograms = scanners.ScanStartupPrograms()
	for _, sp := range results.StartupPrograms {
		if sp.Suspicious {
			results.Flags = append(results.Flags, pkg.Flag{
				Type:     "startup_program",
				Severity: "high",
				Name:     sp.Name,
				Detail:   fmt.Sprintf("Suspicious startup program: %s (%s)", sp.Name, sp.Command),
			})
		}
	}
	fmt.Printf("       Found %d startup entries\n", len(results.StartupPrograms))

	fmt.Println("[6/7] Scanning for suspicious services and installed programs...")
	results.SuspiciousServices = scanners.ScanServices()
	for _, svc := range results.SuspiciousServices {
		results.Flags = append(results.Flags, pkg.Flag{
			Type:     "service",
			Severity: "high",
			Name:     svc,
			Detail:   fmt.Sprintf("Suspicious service: %s", svc),
		})
	}
	fmt.Printf("       Found %d suspicious services\n", len(results.SuspiciousServices))

	results.InstalledPrograms = scanners.ScanInstalledPrograms()
	for _, prog := range results.InstalledPrograms {
		results.Flags = append(results.Flags, pkg.Flag{
			Type:     "installed_program",
			Severity: "high",
			Name:     prog,
			Detail:   fmt.Sprintf("Suspicious installed program: %s", prog),
		})
	}
	fmt.Printf("       Found %d suspicious installed programs\n", len(results.InstalledPrograms))

	fmt.Println("[7/7] Scanning Minecraft mods...")
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

	reportCode := generateReportCode()

	fmt.Println()
	fmt.Println("========================================")
	fmt.Println("Scan complete! Uploading results...")
	fmt.Println("========================================")

	err := pkg.UploadResults("https://swiftac-api.onrender.com", reportCode, results)
	if err != nil {
		fmt.Printf("Failed to upload results: %v\n", err)
		fmt.Println("Results saved locally to scan_result.json")
		saveLocal(results)
		waitAndExit()
		return
	}

	fmt.Println()
	fmt.Println("  ============================================")
	fmt.Println("   ✅  SCAN COMPLETE!")
	fmt.Println("  ============================================")
	fmt.Println()
	fmt.Printf("   Your Report Code:  %s\n", reportCode)
	fmt.Println()
	fmt.Println("   Send this code to the staff member")
	fmt.Println("   who requested the scan.")
	fmt.Println()
	fmt.Println("   They will enter it on the website")
	fmt.Println("   to view your results.")
	fmt.Println()
	fmt.Println("  ============================================")
	fmt.Println()
	waitAndExit()
}

func generateReportCode() string {
	const chars = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	code := make([]byte, 8)
	for i := range code {
		n, _ := rand.Int(rand.Reader, big.NewInt(int64(len(chars))))
		code[i] = chars[n.Int64()]
	}
	return fmt.Sprintf("SWIFT-%s-%s", string(code[:4]), string(code[4:]))
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
