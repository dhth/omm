package cli

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

var (
	errCouldntCreateTempTestDir = errors.New("couldn't create temporary test directory")
	errCouldntBuildBinary       = errors.New("couldn't build binary")
)

func skipIntegration(t *testing.T) {
	t.Helper()
	if os.Getenv("INTEGRATION") != "1" {
		t.Skip("Skipping integration tests")
	}
}

func setUpTestBinary() (string, string, error) {
	var zero string
	tempDir, err := os.MkdirTemp("", "")
	if err != nil {
		return zero, zero, fmt.Errorf("%w: %s", errCouldntCreateTempTestDir, err.Error())
	}

	binPath := filepath.Join(tempDir, "omm")
	buildArgs := []string{"build", "-o", binPath, "../.."}

	c := exec.Command("go", buildArgs...)
	err = c.Run()
	if err != nil {
		return zero, zero, fmt.Errorf("%w: %s", errCouldntBuildBinary, err.Error())
	}

	return tempDir, binPath, nil
}
