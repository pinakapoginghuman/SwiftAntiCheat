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
	"sync"
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

var spinner = []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"}

func main() {
	fmt.Println(Banner)

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

	fmt.Print("  Initializing scan...")
	spin(500)
	fmt.Print("\r                              \r")

	var wg sync.WaitGroup
	var mu sync.Mutex

	wg.Add(5)

	go func() {
		defer wg.Done()
		info := scanners.GetSystemInfo()
		mu.Lock()
		results.SystemInfo = info
		results.HWIDHash = generateHWID(info)
		mu.Unlock()
	}()

	go func() {
		defer wg.Done()
		procs := scanners.ScanProcesses()
		mu.Lock()
		results.Processes = procs
		for _, proc := range procs {
			if proc.Suspicious {
				results.Flags = append(results.Flags, pkg.Flag{
					Type: "suspicious_process", Severity: "high",
					Name: proc.Name, Detail: fmt.Sprintf("Suspicious process: %s", proc.Name),
				})
			}
		}
		mu.Unlock()
	}()

	go func() {
		defer wg.Done()
		files := scanners.ScanFilesystem()
		mu.Lock()
		results.SuspiciousFiles = files
		for _, f := range files {
			results.Flags = append(results.Flags, pkg.Flag{
				Type: "suspicious_file", Severity: f.Severity,
				Name: f.Name, Detail: fmt.Sprintf("Found: %s", f.Path),
			})
		}
		mu.Unlock()
	}()

	go func() {
		defer wg.Done()
		artifacts := scanners.ScanWindowsArtifacts()
		mu.Lock()
		results.WindowsArtifacts = artifacts
		for _, pf := range artifacts.PrefetchFiles {
			if pf.Suspicious {
				results.Flags = append(results.Flags, pkg.Flag{
					Type: "prefetch", Severity: "medium",
					Name: pf.Name, Detail: fmt.Sprintf("Suspicious prefetch: %s", pf.Name),
				})
			}
		}
		for _, rk := range artifacts.SuspiciousRegistryKeys {
			results.Flags = append(results.Flags, pkg.Flag{
				Type: "registry", Severity: "high", Name: rk, Detail: fmt.Sprintf("Registry: %s", rk),
			})
		}
		for _, rk := range artifacts.SuspiciousRunKeys {
			results.Flags = append(results.Flags, pkg.Flag{
				Type: "startup", Severity: "high", Name: rk, Detail: fmt.Sprintf("Startup: %s", rk),
			})
		}
		if artifacts.EventLogCleared {
			results.Flags = append(results.Flags, pkg.Flag{
				Type: "event_log_cleared", Severity: "high",
				Name: "Event Log Cleared", Detail: "Event log was cleared",
			})
		}
		for _, dns := range artifacts.DnsCache {
			results.Flags = append(results.Flags, pkg.Flag{
				Type: "dns_cache", Severity: "medium", Name: dns, Detail: fmt.Sprintf("DNS: %s", dns),
			})
		}
		mu.Unlock()
	}()

	go func() {
		defer wg.Done()
		startupProgs := scanners.ScanStartupPrograms()
		suspiciousServices := scanners.ScanServices()
		installedProgs := scanners.ScanInstalledPrograms()
		mu.Lock()
		results.StartupPrograms = startupProgs
		results.SuspiciousServices = suspiciousServices
		results.InstalledPrograms = installedProgs
		for _, sp := range startupProgs {
			if sp.Suspicious {
				results.Flags = append(results.Flags, pkg.Flag{
					Type: "startup_program", Severity: "high",
					Name: sp.Name, Detail: fmt.Sprintf("Startup: %s", sp.Name),
				})
			}
		}
		for _, svc := range suspiciousServices {
			results.Flags = append(results.Flags, pkg.Flag{
				Type: "service", Severity: "high", Name: svc, Detail: fmt.Sprintf("Service: %s", svc),
			})
		}
		for _, prog := range installedProgs {
			results.Flags = append(results.Flags, pkg.Flag{
				Type: "installed_program", Severity: "high", Name: prog, Detail: fmt.Sprintf("Program: %s", prog),
			})
		}
		mu.Unlock()
	}()

	// Animated spinner while scanning
	done := make(chan bool)
	go func() {
		i := 0
		for {
			select {
			case <-done:
				return
			default:
				fmt.Printf("\r  %s Scanning your system...", spinner[i%len(spinner)])
				i++
				time.Sleep(100 * time.Millisecond)
			}
		}
	}()

	wg.Wait()
	close(done)

	// Now scan Minecraft mods (with visible progress)
	fmt.Print("\r  Scanning Minecraft files...   ")
	results.MinecraftMods = scanners.ScanMinecraftMods()
	for _, mod := range results.MinecraftMods {
		if mod.Suspicious {
			results.Flags = append(results.Flags, pkg.Flag{
				Type: "minecraft_mod", Severity: "high",
				Name: mod.Name, Detail: fmt.Sprintf("Suspicious mod: %s", mod.Name),
			})
		}
	}
	fmt.Printf("\r  ✓ Minecraft: %d mods checked\n", len(results.MinecraftMods))

	reportCode := generateReportCode()

	fmt.Println()
	fmt.Println("  Uploading results...")
	fmt.Print("  ")
	spin(800)

	err := pkg.UploadResults("https://swiftac-api.onrender.com", reportCode, results)
	if err != nil {
		fmt.Printf("\r  ✗ Upload failed: %v\n", err)
		saveLocal(results)
		waitAndExit()
		return
	}

	fmt.Print("\r  ✓ Results uploaded             \n")
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
	fmt.Println("  ============================================")
	fmt.Println()
	waitAndExit()
}

func spin(ms time.Duration) {
	for i := 0; i < 8; i++ {
		fmt.Printf("%s", spinner[i%len(spinner)])
		time.Sleep(ms / 8)
	}
	fmt.Print("\b \b")
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

func saveLocal(results pkg.ScanReport) {
	data, _ := json.MarshalIndent(results, "", "  ")
	os.WriteFile("scan_result.json", data, 0644)
}

func waitAndExit() {
	fmt.Println()
	fmt.Println("  Press Enter to exit...")
	fmt.Scanln()
	os.Exit(0)
}
