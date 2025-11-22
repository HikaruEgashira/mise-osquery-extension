package scanner

import (
	"log"
	"os"
	"path/filepath"
	"sync"
)

type Package struct {
	Tool        string
	Version     string
	Manager     string
	InstallPath string
}

func ScanAllManagers() ([]Package, error) {
	var allPackages []Package
	var mu sync.Mutex
	var wg sync.WaitGroup

	managers := []struct {
		name    string
		scanner func() ([]Package, error)
	}{
		{"mise", ScanMise},
		{"asdf", ScanAsdf},
	}

	// Scan all managers concurrently
	for _, mgr := range managers {
		wg.Add(1)
		go func(name string, scanFunc func() ([]Package, error)) {
			defer wg.Done()
			packages, err := scanFunc()
			if err != nil {
				log.Printf("Error scanning %s: %v", name, err)
				return
			}
			mu.Lock()
			allPackages = append(allPackages, packages...)
			mu.Unlock()
		}(mgr.name, mgr.scanner)
	}

	wg.Wait()
	return allPackages, nil
}

func scanInstallsDirectory(basePath string, manager string) ([]Package, error) {
	packages := []Package{}

	if _, err := os.Stat(basePath); os.IsNotExist(err) {
		return packages, nil
	}

	// Read tool directories (e.g., node, python, ruby)
	toolDirs, err := os.ReadDir(basePath)
	if err != nil {
		return packages, nil
	}

	for _, toolDir := range toolDirs {
		if !toolDir.IsDir() {
			continue
		}

		toolName := toolDir.Name()
		toolPath := filepath.Join(basePath, toolName)

		// Read version directories
		versionDirs, err := os.ReadDir(toolPath)
		if err != nil {
			continue
		}

		for _, versionDir := range versionDirs {
			if !versionDir.IsDir() {
				continue
			}

			version := versionDir.Name()
			installPath := filepath.Join(toolPath, version)

			// Verify the installation directory exists and is valid
			if info, err := os.Stat(installPath); err == nil && info.IsDir() {
				packages = append(packages, Package{
					Tool:        toolName,
					Version:     version,
					Manager:     manager,
					InstallPath: installPath,
				})
			}
		}
	}

	return packages, nil
}

// ScanMise scans mise installations
func ScanMise() ([]Package, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return []Package{}, nil
	}

	// mise stores installations in ~/.local/share/mise/installs
	miseInstallsPath := filepath.Join(home, ".local", "share", "mise", "installs")

	// Also check MISE_DATA_DIR environment variable
	if miseDataDir := os.Getenv("MISE_DATA_DIR"); miseDataDir != "" {
		miseInstallsPath = filepath.Join(miseDataDir, "installs")
	}

	return scanInstallsDirectory(miseInstallsPath, "mise")
}

// ScanAsdf scans asdf installations
func ScanAsdf() ([]Package, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return []Package{}, nil
	}

	// asdf stores installations in ~/.asdf/installs
	asdfInstallsPath := filepath.Join(home, ".asdf", "installs")

	// Also check ASDF_DATA_DIR environment variable
	if asdfDataDir := os.Getenv("ASDF_DATA_DIR"); asdfDataDir != "" {
		asdfInstallsPath = filepath.Join(asdfDataDir, "installs")
	}

	return scanInstallsDirectory(asdfInstallsPath, "asdf")
}
