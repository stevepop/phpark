package php

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

// InstallPHP installs a PHP version with FPM
func InstallPHP(version string) error {
	fmt.Printf("📥 Installing PHP %s-FPM...\n", version)

	packageName := fmt.Sprintf("php%s-fpm", version)

	// Try installing directly from default repos first.
	// Ubuntu 24.04 ships PHP 8.3; this avoids any PPA setup on those systems.
	fmt.Println("   Trying default repositories...")
	cmd := exec.Command("apt-get", "install", "-y", packageName)
	if err := cmd.Run(); err != nil {
		// Not in default repos — add the ondrej/php repository manually.
		// We bypass add-apt-repository (which contacts api.launchpad.net via
		// Python's httplib2) and add the repo directly from packages.sury.org.
		// This is the same maintainer, same packages, no Launchpad API call.
		fmt.Println("   Not in default repos, adding PHP repository...")
		if err := addSuryPHPRepo(); err != nil {
			return fmt.Errorf("failed to add PHP repository: %w", err)
		}

		// Update package list after adding repo
		fmt.Println("   Updating package list...")
		cmd = exec.Command("apt-get", "update")
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("failed to update packages: %w", err)
		}

		// Retry install from the new repo
		fmt.Printf("   Installing %s...\n", packageName)
		cmd = exec.Command("apt-get", "install", "-y", packageName)
		if out, err := cmd.CombinedOutput(); err != nil {
			return fmt.Errorf("failed to install PHP %s: %w\n   %s", version, err, strings.TrimSpace(string(out)))
		}
	}

	// Install common extensions
	fmt.Println("   Installing common extensions...")
	extensions := []string{
		fmt.Sprintf("php%s-cli", version),
		fmt.Sprintf("php%s-common", version),
		fmt.Sprintf("php%s-mysql", version),
		fmt.Sprintf("php%s-curl", version),
		fmt.Sprintf("php%s-mbstring", version),
		fmt.Sprintf("php%s-xml", version),
		fmt.Sprintf("php%s-zip", version),
	}

	for _, ext := range extensions {
		cmd = exec.Command("apt-get", "install", "-y", ext)
		cmd.Run() // Non-fatal if individual extensions fail
	}

	fmt.Printf("\n✅ PHP %s installed successfully!\n", version)
	return nil
}

// addSuryPHPRepo adds the ondrej/php repository directly from packages.sury.org,
// bypassing add-apt-repository which requires a live connection to api.launchpad.net.
// packages.sury.org is maintained by the same author (Ondřej Surý) and contains
// identical packages.
func addSuryPHPRepo() error {
	// Get Ubuntu codename (e.g. "jammy", "noble")
	out, err := exec.Command("lsb_release", "-cs").Output()
	if err != nil {
		return fmt.Errorf("failed to get Ubuntu codename: %w", err)
	}
	codename := strings.TrimSpace(string(out))

	// Ensure gnupg and wget are available for key import
	exec.Command("apt-get", "install", "-y", "--no-install-recommends", "gnupg", "wget").Run()

	// Create keyrings directory
	if err := os.MkdirAll("/etc/apt/keyrings", 0755); err != nil {
		return fmt.Errorf("failed to create keyrings directory: %w", err)
	}

	// Download and store the signing key
	keyCmd := exec.Command("sh", "-c",
		`wget -qO- https://packages.sury.org/php/apt.gpg > /etc/apt/keyrings/sury-php.gpg`)
	if out, err := keyCmd.CombinedOutput(); err != nil {
		return fmt.Errorf("failed to fetch PHP repo signing key: %w\n   %s", err, strings.TrimSpace(string(out)))
	}

	// Write apt sources entry
	source := fmt.Sprintf(
		"deb [signed-by=/etc/apt/keyrings/sury-php.gpg] https://packages.sury.org/php/ %s main\n",
		codename)
	if err := os.WriteFile("/etc/apt/sources.list.d/sury-php.list", []byte(source), 0644); err != nil {
		return fmt.Errorf("failed to write apt source file: %w", err)
	}

	return nil
}

// PromptInstallPHP asks user if they want to install a PHP version
func PromptInstallPHP(version string) (bool, error) {
	fmt.Printf("\n⚠️  PHP %s is not installed.\n", version)
	fmt.Printf("   Would you like to install it now? (y/N): ")

	var response string
	fmt.Scanln(&response)

	return response == "y" || response == "Y" || response == "yes", nil
}
