package scanner

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestScanInstallsDirectory(t *testing.T) {
	// Create a temporary test directory structure
	tmpDir := t.TempDir()
	installsPath := filepath.Join(tmpDir, "installs")

	// Create some test tool installations
	nodePath := filepath.Join(installsPath, "node", "20.10.0")
	pythonPath := filepath.Join(installsPath, "python", "3.12.0")

	err := os.MkdirAll(nodePath, 0755)
	assert.NoError(t, err)
	err = os.MkdirAll(pythonPath, 0755)
	assert.NoError(t, err)

	// Test scanning the directory
	packages, err := scanInstallsDirectory(installsPath, "mise")
	assert.NoError(t, err)
	assert.Len(t, packages, 2)

	// Verify package details
	foundNode := false
	foundPython := false
	for _, pkg := range packages {
		if pkg.Tool == "node" && pkg.Version == "20.10.0" {
			foundNode = true
			assert.Equal(t, "mise", pkg.Manager)
			assert.Equal(t, nodePath, pkg.InstallPath)
		}
		if pkg.Tool == "python" && pkg.Version == "3.12.0" {
			foundPython = true
			assert.Equal(t, "mise", pkg.Manager)
			assert.Equal(t, pythonPath, pkg.InstallPath)
		}
	}
	assert.True(t, foundNode)
	assert.True(t, foundPython)
}

func TestScanNonExistentDirectory(t *testing.T) {
	packages, err := scanInstallsDirectory("/nonexistent/path", "mise")
	assert.NoError(t, err)
	assert.Empty(t, packages)
}

func TestScanMise(t *testing.T) {
	// This test will only work if mise is installed on the system
	// Skip if MISE_DATA_DIR is not set and ~/.local/share/mise doesn't exist
	home, err := os.UserHomeDir()
	if err != nil {
		t.Skip("Cannot get user home directory")
	}

	miseInstallsPath := filepath.Join(home, ".local", "share", "mise", "installs")
	if miseDataDir := os.Getenv("MISE_DATA_DIR"); miseDataDir != "" {
		miseInstallsPath = filepath.Join(miseDataDir, "installs")
	}

	if _, err := os.Stat(miseInstallsPath); os.IsNotExist(err) {
		t.Skip("mise installations not found")
	}

	packages, err := ScanMise()
	assert.NoError(t, err)
	// We can't assert the exact number of packages, but we can verify the structure
	if len(packages) > 0 {
		pkg := packages[0]
		assert.NotEmpty(t, pkg.Tool)
		assert.NotEmpty(t, pkg.Version)
		assert.Equal(t, "mise", pkg.Manager)
		assert.NotEmpty(t, pkg.InstallPath)
	}
}

func TestScanAsdf(t *testing.T) {
	// This test will only work if asdf is installed on the system
	// Skip if ASDF_DATA_DIR is not set and ~/.asdf doesn't exist
	home, err := os.UserHomeDir()
	if err != nil {
		t.Skip("Cannot get user home directory")
	}

	asdfInstallsPath := filepath.Join(home, ".asdf", "installs")
	if asdfDataDir := os.Getenv("ASDF_DATA_DIR"); asdfDataDir != "" {
		asdfInstallsPath = filepath.Join(asdfDataDir, "installs")
	}

	if _, err := os.Stat(asdfInstallsPath); os.IsNotExist(err) {
		t.Skip("asdf installations not found")
	}

	packages, err := ScanAsdf()
	assert.NoError(t, err)
	// We can't assert the exact number of packages, but we can verify the structure
	if len(packages) > 0 {
		pkg := packages[0]
		assert.NotEmpty(t, pkg.Tool)
		assert.NotEmpty(t, pkg.Version)
		assert.Equal(t, "asdf", pkg.Manager)
		assert.NotEmpty(t, pkg.InstallPath)
	}
}

func TestScanAllManagers(t *testing.T) {
	packages, err := ScanAllManagers()
	assert.NoError(t, err)
	// We can't assert the exact number of packages, but we can verify the structure
	for _, pkg := range packages {
		assert.NotEmpty(t, pkg.Tool)
		assert.NotEmpty(t, pkg.Version)
		assert.Contains(t, []string{"mise", "asdf"}, pkg.Manager)
		assert.NotEmpty(t, pkg.InstallPath)
	}
}
