package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// installBrewPackages reads a given file and installs each line as a Homebrew package.
func installBrewPackages(packagesFile string, isCask bool) {
	file, err := os.Open(packagesFile)
	if err != nil {
		log.Printf("Warning: Could not open package list %s: %v. Skipping.", packagesFile, err)
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	installType := "formula"
	args := []string{"install"}
	if isCask {
		installType = "cask"
		args = append(args, "--cask")
	}

	fmt.Printf("
--- Installing Homebrew %ss ---
", installType)

	for scanner.Scan() {
		packageName := strings.TrimSpace(scanner.Text())
		if packageName == "" {
			continue
		}

		fmt.Printf("
Installing %s: %s
", installType, packageName)
		cmdArgs := append(args, packageName)
		cmd := exec.Command("brew", cmdArgs...)
		cmd.Stdout = os.Stdout // Show brew's output directly
		cmd.Stderr = os.Stderr // Show brew's errors directly

		if err := cmd.Run(); err != nil {
			log.Printf("Error installing %s. It may already be installed or there was another issue.", packageName)
		}
	}

	if err := scanner.Err(); err != nil {
		log.Printf("Error reading package list %s: %v", packagesFile, err)
	}
}

// linkDotfiles discovers and symlinks configuration directories.
func linkDotfiles(dotfilesDir, configDestDir string) {
	configSourceDir := filepath.Join(dotfilesDir, "config")

	fmt.Printf("
--- Starting Dotfile Symlinking ---
")
	fmt.Printf("Reading configurations from: %s
", configSourceDir)

	entries, err := os.ReadDir(configSourceDir)
	if err != nil {
		log.Fatalf("Error: Could not read config source directory %s. %v", configSourceDir, err)
	}

	if err := os.MkdirAll(configDestDir, 0755); err != nil {
		log.Fatalf("Error: Could not create destination ~/.config directory. %v", err)
	}

	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}
		configName := entry.Name()
		sourcePath := filepath.Join(configSourceDir, configName)
		destPath := filepath.Join(configDestDir, configName)

		fmt.Printf("Processing %s...
", configName)

		if _, err := os.Lstat(destPath); err == nil {
			fmt.Printf("  - Removing existing target: %s
", destPath)
			if err := os.RemoveAll(destPath); err != nil {
				log.Printf("  - Warning: Failed to remove existing destination %s. %v
", destPath, err)
				continue
			}
		}

		fmt.Printf("  - Linking %s -> %s
", sourcePath, destPath)
		if err := os.Symlink(sourcePath, destPath); err != nil {
			log.Printf("  - Error: Failed to link %s. %v
", configName, err)
		} else {
			fmt.Printf("  - Successfully linked %s.
", configName)
		}
	}
}

func main() {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Fatalf("Error: Could not get user home directory. %v", err)
	}

	dotfilesDir := filepath.Join(homeDir, "Workspace", "dotfiles")
	configDestDir := filepath.Join(homeDir, ".config")

	// Step 1: Install Homebrew packages
	if _, err := exec.LookPath("brew"); err != nil {
		log.Println("Warning: Homebrew command 'brew' not found in PATH. Skipping all package installations.")
		log.Println("Please install Homebrew first by visiting https://brew.sh/")
	} else {
		installBrewPackages(filepath.Join(dotfilesDir, "brew_formulae.txt"), false)
		installBrewPackages(filepath.Join(dotfilesDir, "brew_casks.txt"), true)
	}

	// Step 2: Link the dotfiles
	linkDotfiles(dotfilesDir, configDestDir)

	fmt.Println("
--- Restore Script Complete ---")
}
