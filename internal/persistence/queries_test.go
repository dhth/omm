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

var testDB *sql.DB

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
	for _, tbl := range []string{"task"} {
		_, err = testDB.Exec(fmt.Sprintf("DELETE FROM %s", tbl))
		if err != nil {
			t.Fatalf("failed to clean up table %q: %v", tbl, err)
		}
		_, err := testDB.Exec("DELETE FROM sqlite_sequence WHERE name=?;", tbl)
		if err != nil {
			t.Fatalf("failed to reset auto increment for table %q: %v", tbl, err)
		}
	}
	_, err = testDB.Exec(`UPDATE task_sequence
SET sequence = '[]'
WHERE id = 1;`)
	if err != nil {
		t.Fatalf("failed to clean up table task_sequence: %v", err)
	}
}

func getSampleTasks() ([]types.Task, int, int) {
	numActive := 3
	numInactive := 2

	tasks := make([]types.Task, numActive+numInactive)
	contexts := make([]string, numActive+numInactive)
	now := time.Now().UTC()
	counter := 0
	for range numActive {
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
	for range numInactive {
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

	return tasks, numActive, numInactive
}

func seedDB(t *testing.T, db *sql.DB) (int, int) {
	t.Helper()

	tasks, na, ni := getSampleTasks()

	for _, task := range tasks {
		_, err := db.Exec(`
INSERT INTO task (summary, active, created_at, updated_at)
VALUES (?, ?, ?, ?)`, task.Summary, task.Active, task.CreatedAt, task.UpdatedAt)
		if err != nil {
			t.Fatalf("failed to insert data into table \"task\": %v", err)
		}
	}

	seqItems := make([]int, na)
	for i := range na {
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

	return na, ni
}

func TestInsertTasksWorksWithEmptyTaskList(t *testing.T) {
	t.Cleanup(func() { cleanupDB(t) })

	// GIVEN
	// WHEN
	now := time.Now().UTC()
	tasks := []types.Task{
		{
			Summary:   "prefix: new task 1",
			Active:    true,
			CreatedAt: now,
			UpdatedAt: now,
		},
		{
			Summary:   "prefix: new inactive task 1",
			Active:    false,
			CreatedAt: now,
			UpdatedAt: now,
		},
		{
			Summary:   "prefix: new task 3",
			Active:    true,
			CreatedAt: now,
			UpdatedAt: now,
		},
	}
	lastID, err := InsertTasks(testDB, tasks, true)
	assert.Equal(t, lastID, int64(3), "last ID is not correct")
	require.NoError(t, err)

	// THEN
	numActiveRes, err := fetchNumActiveTasks(testDB)
	require.NoError(t, err)
	assert.Equal(t, numActiveRes, 2, "number of active tasks didn't increase by the correct amount")

	numTotalRes, err := fetchNumTotalTasks(testDB)
	require.NoError(t, err)
	assert.Equal(t, numTotalRes, 3, "number of total tasks didn't increase by the correct amount")

	lastTask, err := fetchTaskByID(testDB, lastID)
	require.NoError(t, err)
	assert.Equal(t, tasks[2].Active, lastTask.Active)
	assert.Equal(t, tasks[2].Summary, lastTask.Summary)
	assert.Equal(t, tasks[2].Context, lastTask.Context)

	seq, err := fetchTaskSequence(testDB)
	require.NoError(t, err)
	assert.Equal(t, seq, []uint64{1, 3}, "task sequence isn't correct")
}

func TestInsertTasksAddsTasksAtTheTop(t *testing.T) {
	t.Cleanup(func() { cleanupDB(t) })

	// GIVEN
	na, ni := seedDB(t, testDB)

	// WHEN
	now := time.Now().UTC()
	tasks := []types.Task{
		{
			Summary:   "prefix: new task 1",
			Active:    true,
			CreatedAt: now,
			UpdatedAt: now,
		},
		{
			Summary:   "prefix: new inactive task 1",
			Active:    false,
			CreatedAt: now,
			UpdatedAt: now,
		},
		{
			Summary:   "prefix: new task 3",
			Active:    true,
			CreatedAt: now,
			UpdatedAt: now,
		},
	}

	_, err := InsertTasks(testDB, tasks, true)
	require.NoError(t, err)

	// THEN
	numActiveRes, err := fetchNumActiveTasks(testDB)
	require.NoError(t, err)
	assert.Equal(t, numActiveRes, na+2, "number of active tasks didn't increase by the correct amount")

	numTotalRes, err := fetchNumTotalTasks(testDB)
	require.NoError(t, err)
	assert.Equal(t, numTotalRes, na+ni+3, "number of total tasks didn't increase by the correct amount")

	seq, err := fetchTaskSequence(testDB)
	require.NoError(t, err)
	assert.Equal(t, seq, []uint64{6, 8, 1, 2, 3}, "task sequence isn't correct")
}

func TestInsertTasksAddsTasksAtTheEnd(t *testing.T) {
	t.Cleanup(func() { cleanupDB(t) })

	// GIVEN
	na, _ := seedDB(t, testDB)

	// WHEN
	now := time.Now().UTC()
	tasks := []types.Task{
		{
			Summary:   "prefix: new task 1",
			Active:    true,
			CreatedAt: now,
			UpdatedAt: now,
		},
		{
			Summary:   "prefix: new task 2",
			Active:    true,
			CreatedAt: now,
			UpdatedAt: now,
		},
	}

	_, err := InsertTasks(testDB, tasks, false)
	require.NoError(t, err)

	// THEN
	numActiveRes, err := fetchNumActiveTasks(testDB)
	require.NoError(t, err)
	assert.Equal(t, numActiveRes, na+2, "number of active tasks didn't increase by the correct amount")

	seq, err := fetchTaskSequence(testDB)
	require.NoError(t, err)
	assert.Equal(t, seq, []uint64{1, 2, 3, 6, 7}, "task sequence isn't correct")
}
