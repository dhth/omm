package cli

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/google/uuid"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTasksCmd(t *testing.T) {
	skipIntegration(t)

	tempDir, binPath, err := setUpTestBinary()
	if err != nil {
		require.NoErrorf(t, err, "error setting up test binary: %w", err)
	}

	defer func() {
		err := os.RemoveAll(tempDir)
		if err != nil {
			t.Fatalf("couldn't clean up temporary directory (%s): %s", binPath, err)
		}
	}()

	//-------------//
	//  SUCCESSES  //
	//-------------//

	t.Run("Help", func(t *testing.T) {
		// GIVEN
		// WHEN
		c := exec.Command(binPath, "tasks", "-h")
		err := c.Run()

		// THEN
		assert.NoError(t, err)
	})

	t.Run("Listing tasks works", func(t *testing.T) {
		// GIVEN
		dbPath := filepath.Join(tempDir, fmt.Sprintf("omm-%s.db", uuid.New().String()))

		numTasks := 10
		for i := range numTasks {
			c := exec.Command(binPath, "-d", dbPath, fmt.Sprintf("prefix: task %d", 9-i))
			err := c.Run()
			require.NoError(t, err)
		}

		// WHEN
		c := exec.Command(binPath, "-d", dbPath, "tasks", "list")
		output, err := c.CombinedOutput()

		// THEN
		require.NoError(t, err)

		numLines := len(strings.Split(strings.TrimSpace(string(output)), "\n"))
		assert.Equal(t, numTasks, numLines)
	})

	t.Run("Showing task details for a valid index works", func(t *testing.T) {
		// GIVEN
		dbPath := filepath.Join(tempDir, fmt.Sprintf("omm-%s.db", uuid.New().String()))

		numTasks := 10
		for i := range numTasks {
			c := exec.Command(binPath, "-d", dbPath, fmt.Sprintf("prefix: task %d", 10-i))
			err := c.Run()
			require.NoError(t, err)
		}

		// WHEN
		c := exec.Command(binPath, "-d", dbPath, "tasks", "show", "2", "-f", "plain")
		output, err := c.CombinedOutput()

		// THEN
		require.NoError(t, err)

		assert.Equal(t, "prefix: task 2", strings.TrimSpace(string(output)))
	})

	t.Run("Showing task details in json format works", func(t *testing.T) {
		// GIVEN
		dbPath := filepath.Join(tempDir, fmt.Sprintf("omm-%s.db", uuid.New().String()))

		numTasks := 10
		for i := range numTasks {
			c := exec.Command(binPath, "-d", dbPath, fmt.Sprintf("prefix: task %d", 10-i))
			err := c.Run()
			require.NoError(t, err)
		}

		// WHEN
		c := exec.Command(binPath, "-d", dbPath, "tasks", "show", "2", "-f", "json")
		output, err := c.CombinedOutput()

		// THEN
		require.NoError(t, err)

		expected := `
{
  "summary": "prefix: task 2",
  "context": null
}
`
		assert.Equal(t, strings.TrimSpace(expected), strings.TrimSpace(string(output)))
	})

	//------------//
	//  FAILURES  //
	//------------//

	t.Run("Showing task details for index 0 fails", func(t *testing.T) {
		// GIVEN
		dbPath := filepath.Join(tempDir, fmt.Sprintf("omm-%s.db", uuid.New().String()))

		// WHEN
		c := exec.Command(binPath, "-d", dbPath, "tasks", "show", "0")
		output, err := c.CombinedOutput()

		// THEN
		assert.Error(t, err)

		assert.Contains(t, string(output), "invalid value for task index provided;")
	})

	t.Run("Showing task details for an invalid index fails", func(t *testing.T) {
		// GIVEN
		dbPath := filepath.Join(tempDir, fmt.Sprintf("omm-%s.db", uuid.New().String()))

		numTasks := 10
		for i := range numTasks {
			c := exec.Command(binPath, "-d", dbPath, fmt.Sprintf("prefix: task %d", 10-i))
			err := c.Run()
			require.NoError(t, err)
		}

		// WHEN
		c := exec.Command(binPath, "-d", dbPath, "tasks", "show", "11")
		output, err := c.CombinedOutput()

		// THEN
		assert.Error(t, err)

		assert.Contains(t, string(output), "no task exists at given index")
	})

	t.Run("Showing task details with an invalid output format fails", func(t *testing.T) {
		// GIVEN
		dbPath := filepath.Join(tempDir, fmt.Sprintf("omm-%s.db", uuid.New().String()))

		// WHEN
		c := exec.Command(binPath, "-d", dbPath, "tasks", "show", "1", "-f", "unknown")
		output, err := c.CombinedOutput()

		// THEN
		assert.Error(t, err)

		assert.Contains(t, string(output), "invalid value for output format provided;")
	})
}
