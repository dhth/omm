package persistence

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/dhth/omm/internal/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	_ "modernc.org/sqlite" // sqlite driver
)

var (
	testDB          *sql.DB
	numSeedActive   = 3
	numSeedInActive = 2
)

func TestMain(m *testing.M) {
	var err error
	testDB, err = sql.Open("sqlite", ":memory:")
	if err != nil {
		panic(err)
	}

	err = InitDB(testDB)
	if err != nil {
		panic(err)
	}
	err = UpgradeDB(testDB, 1)
	if err != nil {
		panic(err)
	}
	code := m.Run()

	testDB.Close()

	os.Exit(code)
}

func cleanupDB(t *testing.T) {
	var err error
	for _, tbl := range []string{"task", "task_sequence"} {
		_, err = testDB.Exec(fmt.Sprintf("DELETE FROM %s", tbl))
		if err != nil {
			t.Fatalf("failed to clean up table %q: %v", tbl, err)
		}
	}
}

func seedDB(t *testing.T, db *sql.DB) {
	t.Helper()

	tasks := make([]types.Task, numSeedActive+numSeedInActive)
	contexts := make([]string, numSeedActive+numSeedInActive)
	now := time.Now().UTC()
	counter := 0
	for range numSeedActive {
		contexts[counter] = fmt.Sprintf("context for task %d", counter)
		tasks[counter] = types.Task{
			Summary:   fmt.Sprintf("prefix: task %d", counter),
			Active:    true,
			Context:   &contexts[counter],
			CreatedAt: now,
			UpdatedAt: now,
		}
		counter++
	}
	for range numSeedInActive {
		contexts[counter] = fmt.Sprintf("context for task %d", counter)
		tasks[counter] = types.Task{
			Summary:   fmt.Sprintf("prefix: task %d", counter),
			Active:    false,
			Context:   &contexts[counter],
			CreatedAt: now,
			UpdatedAt: now,
		}
		counter++
	}
	for _, task := range tasks {
		_, err := db.Exec(`
INSERT INTO task (summary, active, created_at, updated_at)
VALUES (?, ?, ?, ?)`, task.Summary, task.Active, task.CreatedAt, task.UpdatedAt)
		if err != nil {
			t.Fatalf("failed to insert data into table \"task\": %v", err)
		}
	}

	seqItems := make([]int, numSeedActive)
	for i := range numSeedActive {
		seqItems[i] = i + 1
	}
	sequenceJSON, err := json.Marshal(seqItems)
	if err != nil {
		t.Fatalf("failed to marshall JSON data for seeding: %v", err)
	}

	_, err = db.Exec(`
UPDATE task_sequence
SET sequence = ?
WHERE id = 1;
`, sequenceJSON)
	if err != nil {
		t.Fatalf("failed to insert data into table \"task_sequence\": %v", err)
	}
}

func TestImportTask(t *testing.T) {
	t.Cleanup(func() { cleanupDB(t) })

	// GIVEN
	seedDB(t, testDB)
	numActiveTasksBefore, err := fetchNumActiveTasks(testDB)
	require.NoError(t, err)
	numTotalTasksBefore, err := fetchNumTotalTasks(testDB)
	require.NoError(t, err)

	// WHEN
	summary := "prefix: an imported task"
	now := time.Now().UTC()
	err = ImportTask(testDB, summary, true, now, now)
	require.NoError(t, err)

	// THEN
	numActiveTasksAfter, err := fetchNumActiveTasks(testDB)
	require.NoError(t, err)
	assert.Equal(t, numActiveTasksAfter, numActiveTasksBefore+1, "number of active tasks didn't increase by 1")

	task, err := fetchTaskByID(testDB, numTotalTasksBefore+1)
	require.NoError(t, err)
	assert.True(t, task.Active)
	assert.Equal(t, summary, task.Summary)

	seq, err := fetchTaskSequence(testDB)
	require.NoError(t, err)
	require.Equal(t, numActiveTasksAfter, len(seq), "number of tasks in task sequence doesn't match number of active tasks")
	assert.Equal(t, seq[0], task.ID, "newly added task is not shown at the top of the list")
}

func TestImportTaskSummaries(t *testing.T) {
	t.Cleanup(func() { cleanupDB(t) })

	// GIVEN
	seedDB(t, testDB)
	numActiveTasksBefore, err := fetchNumActiveTasks(testDB)
	require.NoError(t, err)
	numTotalTasksBefore, err := fetchNumTotalTasks(testDB)
	require.NoError(t, err)

	// WHEN
	summaries := []string{
		"prefix: imported task 1",
		"prefix: imported task 2",
		"prefix: imported task 3",
	}
	now := time.Now().UTC()
	err = ImportTaskSummaries(testDB, summaries, true, now, now)
	require.NoError(t, err)

	// THEN
	numActiveTasksAfter, err := fetchNumActiveTasks(testDB)
	require.NoError(t, err)
	assert.Equal(t, numActiveTasksAfter, numActiveTasksBefore+len(summaries), "number of active tasks didn't increase by the correct amount")

	task, err := fetchTaskByID(testDB, numTotalTasksBefore+1)
	require.NoError(t, err)
	assert.True(t, task.Active)
	assert.Equal(t, summaries[0], task.Summary)

	seq, err := fetchTaskSequence(testDB)
	require.NoError(t, err)
	require.Equal(t, numActiveTasksAfter, len(seq), "number of tasks in task sequence doesn't match number of active tasks")
	for i := range summaries {
		assert.Equal(t, numTotalTasksBefore+i+1, int(seq[i]), "task at sequence position %d is incorrect", i+1)
	}
}
