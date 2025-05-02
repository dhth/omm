package persistence

import (
	"database/sql"
	"encoding/json"
	"time"

	"github.com/dhth/omm/internal/types"
)

const (
	TaskNumLimit    = 300
	ContextMaxBytes = 1024 * 1024
)

func fetchTaskSequence(db *sql.DB) ([]uint64, error) {
	var seq []byte
	seqRow := db.QueryRow("SELECT sequence from task_sequence where id=1;")

	err := seqRow.Scan(&seq)
	if err != nil {
		return nil, err
	}

	var seqItems []uint64
	err = json.Unmarshal(seq, &seqItems)
	if err != nil {
		return nil, err
	}
	return seqItems, nil
}

func fetchNumActiveTasks(db *sql.DB) (int, error) {
	var rowCount int
	err := db.QueryRow("SELECT count(*) from task where active is true").Scan(&rowCount)
	return rowCount, err
}

func fetchNumTotalTasks(db *sql.DB) (int, error) {
	var rowCount int
	err := db.QueryRow("SELECT count(*) from task").Scan(&rowCount)
	return rowCount, err
}

func fetchTaskByID(db *sql.DB, ID int64) (types.Task, error) {
	var entry types.Task
	row := db.QueryRow(`
SELECT id, summary, active, context, created_at, updated_at
from task
WHERE id=?;
`, ID)
	err := row.Scan(&entry.ID,
		&entry.Summary,
		&entry.Active,
		&entry.Context,
		&entry.CreatedAt,
		&entry.UpdatedAt,
	)
	return entry, err
}

func FetchNumActiveTasksShown(db *sql.DB) (int, error) {
	row := db.QueryRow(`
SELECT json_array_length(sequence) AS num_tasks
FROM task_sequence where id=1;
`)

	var numTasks int
	err := row.Scan(&numTasks)
	if err != nil {
		return -1, err
	}

	return numTasks, nil
}

func UpdateTaskSequence(db *sql.DB, sequence []uint64) error {
	sequenceJSON, err := json.Marshal(sequence)
	if err != nil {
		return err
	}

	stmt, err := db.Prepare(`
UPDATE task_sequence
SET sequence = ?
WHERE id = 1;
`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(sequenceJSON)
	if err != nil {
		return err
	}

	return nil
}

func InsertTask(db *sql.DB, summary string, createdAt, updatedAt time.Time) (uint64, error) {
	stmt, err := db.Prepare(`
INSERT INTO task (summary, active, created_at, updated_at)
VALUES (?, true, ?, ?);
`)
	if err != nil {
		return 0, err
	}
	defer stmt.Close()

	res, err := stmt.Exec(summary, createdAt.UTC(), updatedAt.UTC())
	if err != nil {
		return 0, err
	}

	li, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}

	return uint64(li), nil
}

func InsertTasks(db *sql.DB, tasks []types.Task, insertAtTop bool) (int64, error) {
	tx, err := db.Begin()
	if err != nil {
		return -1, err
	}
	defer func() {
		_ = tx.Rollback()
	}()

	query := `INSERT INTO task (summary, context, active, created_at, updated_at)
VALUES `

	values := make([]interface{}, 0, len(tasks)*4)

	for i, t := range tasks {
		if i > 0 {
			query += ","
		}
		query += "(?, ?, ?, ?, ?)"
		values = append(values, t.Summary, t.Context, t.Active, t.CreatedAt.UTC(), t.UpdatedAt.UTC())
	}

	query += ";"

	res, err := tx.Exec(query, values...)
	if err != nil {
		return -1, err
	}

	lastInsertID, err := res.LastInsertId()
	if err != nil {
		return -1, err
	}

	var seq []byte
	seqRow := tx.QueryRow("SELECT sequence from task_sequence where id=1;")

	err = seqRow.Scan(&seq)
	if err != nil {
		return -1, err
	}

	var seqItems []int
	err = json.Unmarshal(seq, &seqItems)
	if err != nil {
		return -1, err
	}

	var newTaskIDs []int
	taskID := int(lastInsertID) - len(tasks) + 1
	for _, t := range tasks {
		if t.Active {
			newTaskIDs = append(newTaskIDs, taskID)
		}
		taskID++
	}

	var updatedSeqItems []int
	if insertAtTop {
		updatedSeqItems = append(newTaskIDs, seqItems...)
	} else {
		updatedSeqItems = append(seqItems, newTaskIDs...)
	}

	sequenceJSON, err := json.Marshal(updatedSeqItems)
	if err != nil {
		return -1, err
	}

	seqUpdateStmt, err := tx.Prepare(`
UPDATE task_sequence
SET sequence = ?
WHERE id = 1;
`)
	if err != nil {
		return -1, err
	}
	defer seqUpdateStmt.Close()

	_, err = seqUpdateStmt.Exec(sequenceJSON)
	if err != nil {
		return -1, err
	}

	err = tx.Commit()
	if err != nil {
		return -1, err
	}
	return lastInsertID, nil
}

func UpdateTaskSummary(db *sql.DB, id uint64, summary string, updatedAt time.Time) error {
	stmt, err := db.Prepare(`
UPDATE task
SET summary = ?,
    updated_at = ?
WHERE id = ?
`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(summary, updatedAt.UTC(), id)
	if err != nil {
		return err
	}
	return nil
}

func UpdateTaskContext(db *sql.DB, id uint64, context string, updatedAt time.Time) error {
	stmt, err := db.Prepare(`
UPDATE task
SET context = ?,
    updated_at = ?
WHERE id = ?
`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(context, updatedAt.UTC(), id)
	if err != nil {
		return err
	}
	return nil
}

func UnsetTaskContext(db *sql.DB, id uint64, updatedAt time.Time) error {
	stmt, err := db.Prepare(`
UPDATE task
SET context = NULL,
    updated_at = ?
WHERE id = ?
`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(updatedAt.UTC(), id)
	if err != nil {
		return err
	}
	return nil
}

func ChangeTaskStatus(db *sql.DB, id uint64, active bool, updatedAt time.Time) error {
	stmt, err := db.Prepare(`
UPDATE task
SET active = ?,
    updated_at = ?
WHERE id = ?
`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(active, updatedAt.UTC(), id)
	if err != nil {
		return err
	}
	return nil
}

func FetchActiveTasks(db *sql.DB, limit int) ([]types.Task, error) {
	var tasks []types.Task

	rows, err := db.Query(`
SELECT t.id, t.summary, t.context, t.created_at, t.updated_at
FROM task_sequence s
JOIN json_each(s.sequence) j ON CAST(j.value AS INTEGER) = t.id
JOIN task t ON t.id = j.value
ORDER BY j.key
LIMIT ?;
`, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var entry types.Task
		err = rows.Scan(&entry.ID,
			&entry.Summary,
			&entry.Context,
			&entry.CreatedAt,
			&entry.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		entry.CreatedAt = entry.CreatedAt.Local()
		entry.UpdatedAt = entry.UpdatedAt.Local()
		entry.Active = true
		tasks = append(tasks, entry)

	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return tasks, nil
}

func FetchInActiveTasks(db *sql.DB, limit int) ([]types.Task, error) {
	var tasks []types.Task

	rows, err := db.Query(`
SELECT id, summary, context, created_at, updated_at
FROM task where active is false
ORDER BY updated_at DESC
LIMIT ?;
`, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var entry types.Task
		err = rows.Scan(&entry.ID,
			&entry.Summary,
			&entry.Context,
			&entry.CreatedAt,
			&entry.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		entry.CreatedAt = entry.CreatedAt.Local()
		entry.UpdatedAt = entry.UpdatedAt.Local()
		entry.Active = false
		tasks = append(tasks, entry)

	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return tasks, nil
}

func DeleteTask(db *sql.DB, id uint64) error {
	stmt, err := db.Prepare(`
DELETE from task
WHERE id=?;
`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(id)
	if err != nil {
		return err
	}
	return nil
}
