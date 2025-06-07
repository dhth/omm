package cli

import (
	"fmt"
	"os"
	"os/exec"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRootCmd(t *testing.T) {
	skipIntegration(t)

	tempDir, binPath, err := setUpTestBinary()
	if err != nil {
		require.NoErrorf(t, err, "error setting up test binary: %w", err)
	}

	defer func() {
		err := os.RemoveAll(tempDir)
		if err != nil {
			fmt.Printf("couldn't clean up temporary directory (%s): %s", binPath, err)
		}
	}()

	// SUCCESSES
	t.Run("Help", func(t *testing.T) {
		// GIVEN
		// WHEN
		c := exec.Command(binPath, "-h")
		err := c.Run()

		// THEN
		assert.NoError(t, err)
	})
}
