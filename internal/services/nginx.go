package services

import (
	"fmt"
	"os/exec"
)

// ReloadNginx reloads nginx service
func ReloadNginx() error {
	cmd := exec.Command("sudo", "systemctl", "reload", "nginx")
	if err := cmd.Run(); err != nil {
		// Try alternative reload method
		cmd = exec.Command("sudo", "nginx", "-s", "reload")
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("failed to reload nginx: %w", err)
		}
	}
	return nil
}

// StartNginx starts nginx if not running
func StartNginx() error {
	// Check if running (no sudo needed for status check)
	cmd := exec.Command("systemctl", "is-active", "nginx")
	if err := cmd.Run(); err == nil {
		return nil // Already running
	}

	// Start nginx
	cmd = exec.Command("sudo", "systemctl", "start", "nginx")
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to start nginx: %w", err)
	}

	// Enable on boot
	cmd = exec.Command("sudo", "systemctl", "enable", "nginx")
	cmd.Run() // Non-fatal

	return nil
}
