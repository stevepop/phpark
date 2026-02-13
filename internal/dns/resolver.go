package dns

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

// SetupDNS configures DNS resolution for .test domains
func SetupDNS(domain string) error {
	if runtime.GOOS == "linux" {
		return setupLinuxDNS(domain)
	}
	return setupMacDNS(domain)
}

// RemoveDNS removes DNS configuration for .test domains
func RemoveDNS(domain string) error {
	if runtime.GOOS == "linux" {
		return removeLinuxDNS(domain)
	}
	return removeMacDNS(domain)
}

// CheckDNS verifies if DNS is configured
func CheckDNS(domain string) (bool, error) {
	if runtime.GOOS == "linux" {
		return checkLinuxDNS(domain)
	}
	return checkMacDNS(domain)
}

// === Linux DNS Setup (dnsmasq) ===

func setupLinuxDNS(domain string) error {
	// Check if dnsmasq is installed
	if _, err := exec.LookPath("dnsmasq"); err != nil {
		return fmt.Errorf("dnsmasq not installed. Install with: sudo apt install dnsmasq")
	}

	// Create dnsmasq config
	configPath := fmt.Sprintf("/etc/dnsmasq.d/%s", domain)
	content := fmt.Sprintf("address=/.%s/127.0.0.1\n", domain)

	// Write config (requires sudo)
	cmd := exec.Command("sudo", "tee", configPath)
	cmd.Stdin = strings.NewReader(content)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to create dnsmasq config: %w", err)
	}

	// Restart dnsmasq
	if err := exec.Command("sudo", "systemctl", "restart", "dnsmasq").Run(); err != nil {
		return fmt.Errorf("failed to restart dnsmasq: %w", err)
	}

	return nil
}

func removeLinuxDNS(domain string) error {
	configPath := fmt.Sprintf("/etc/dnsmasq.d/%s", domain)

	if err := exec.Command("sudo", "rm", "-f", configPath).Run(); err != nil {
		return fmt.Errorf("failed to remove dnsmasq config: %w", err)
	}

	// Restart dnsmasq if it's running
	exec.Command("sudo", "systemctl", "restart", "dnsmasq").Run()

	return nil
}

func checkLinuxDNS(domain string) (bool, error) {
	configPath := fmt.Sprintf("/etc/dnsmasq.d/%s", domain)
	_, err := os.Stat(configPath)
	return err == nil, nil
}

// === macOS DNS Setup (resolver) ===
func setupMacDNS(domain string) error {
	// macOS uses /etc/resolver/ for DNS configuration
	resolverDir := "/etc/resolver"
	resolverFile := filepath.Join(resolverDir, domain)

	// Create resolver directory if it doesn't exist
	if err := exec.Command("sudo", "mkdir", "-p", resolverDir).Run(); err != nil {
		return fmt.Errorf("failed to create resolver directory: %w", err)
	}

	// Create resolver config pointing to localhost on port 53
	// Note: This requires dnsmasq or similar to be running on 127.0.0.1:53
	content := "nameserver 127.0.0.1\nport 53\n"
	cmd := exec.Command("sudo", "tee", resolverFile)
	cmd.Stdin = strings.NewReader(content)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to create resolver config: %w", err)
	}

	// Flush DNS cache
	exec.Command("sudo", "dscacheutil", "-flushcache").Run()
	exec.Command("sudo", "killall", "-HUP", "mDNSResponder").Run()

	return nil
}

func removeMacDNS(domain string) error {
	resolverFile := filepath.Join("/etc/resolver", domain)

	if err := exec.Command("sudo", "rm", "-f", resolverFile).Run(); err != nil {
		return fmt.Errorf("failed to remove resolver config: %w", err)
	}

	// Flush DNS cache
	exec.Command("sudo", "dscacheutil", "-flushcache").Run()
	exec.Command("sudo", "killall", "-HUP", "mDNSResponder").Run()

	return nil
}

func checkMacDNS(domain string) (bool, error) {
	resolverFile := filepath.Join("/etc/resolver", domain)
	_, err := os.Stat(resolverFile)
	return err == nil, nil
}

// TestDNSResolution tests if a domain resolves correctly
func TestDNSResolution(hostname string) (bool, error) {
	// Use nslookup to test
	cmd := exec.Command("nslookup", hostname)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return false, nil // Domain doesn't resolve
	}

	// Check if it resolves to 127.0.0.1
	outputStr := string(output)
	return strings.Contains(outputStr, "127.0.0.1"), nil
}
