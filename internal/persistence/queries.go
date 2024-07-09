package persistence

import (
	"database/sql"
	"encoding/json"
	"time"

	"github.com/dhth/omm/internal/types"
)

const (
	TaskNumLimit = 300
)

func FetchNumActiveTasksFromDB(db *sql.DB) (int, error) {

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

func UpdateTaskSequenceInDB(db *sql.DB, sequence []uint64) error {

	sequenceJson, err := json.Marshal(sequence)
	if err != nil {
		return err
	}

	stmt, err := db.Prepare(`
UPDATE task_sequence SET sequence = ? where id = 1;
`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(sequenceJson)

	if err != nil {
		return err
	}

	return nil
}

func InsertTaskInDB(db *sql.DB, summary string, createdAt, updatedAt time.Time) (uint64, error) {

	stmt, err := db.Prepare(`
INSERT into task (summary, active, created_at, updated_at)
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

func ImportTaskIntoDB(db *sql.DB, summary string, active bool, createdAt, updatedAt time.Time) error {

	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer func() {
		_ = tx.Rollback()
	}()

	query := `INSERT INTO task (summary, active, created_at, updated_at)
VALUES (?, ?, ?, ?);`

	res, err := tx.Exec(query, summary, active, createdAt.UTC(), updatedAt.UTC())
	if err != nil {
		return err
	}

	lastInsertId, err := res.LastInsertId()
	if err != nil {
		return err
	}

	var seq []byte
	seqRow := tx.QueryRow("SELECT sequence from task_sequence where id=1;")

	err = seqRow.Scan(&seq)
	if err != nil {
		return err
	}

	var seqItems []int
	err = json.Unmarshal(seq, &seqItems)
	if err != nil {
		return err
	}

	newTaskID := make([]int, 1)
	newTaskID[0] = int(lastInsertId)
	updatedSeqItems := append(newTaskID, seqItems...)

	sequenceJson, err := json.Marshal(updatedSeqItems)
	if err != nil {
		return err
	}

	seqUpdateStmt, err := tx.Prepare(`
UPDATE task_sequence SET sequence = ? where id = 1;
`)
	if err != nil {
		return err
	}
	defer seqUpdateStmt.Close()

	_, err = seqUpdateStmt.Exec(sequenceJson)

	if err != nil {
		return err
	}

	err = tx.Commit()
	if err != nil {
		return err
	}
	return nil
}

func ImportTasksIntoDB(db *sql.DB, tasks []string, active bool, createdAt, updatedAt time.Time) error {

	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer func() {
		_ = tx.Rollback()
	}()

	query := `INSERT INTO task (summary, active, created_at, updated_at)
VALUES `

	values := make([]interface{}, 0, len(tasks)*4)

	ca := createdAt.UTC()
	ua := updatedAt.UTC()

	for i, ts := range tasks {
		if i > 0 {
			query += ","
		}
		query += "(?, ?, ?, ?)"
		values = append(values, ts, active, ca, ua)
	}
	query += ";"

	res, err := tx.Exec(query, values...)
	if err != nil {
		return err
	}

	lastInsertId, err := res.LastInsertId()
	if err != nil {
		return err
	}

	var seq []byte
	seqRow := tx.QueryRow("SELECT sequence from task_sequence where id=1;")

	err = seqRow.Scan(&seq)
	if err != nil {
		return err
	}

	var seqItems []int
	err = json.Unmarshal(seq, &seqItems)
	if err != nil {
		return err
	}

	newTaskIDs := make([]int, len(tasks))
	counter := 0
	for i := int(lastInsertId) - len(tasks) + 1; i <= int(lastInsertId); i++ {
		newTaskIDs[counter] = i
		counter++
	}
	updatedSeqItems := append(newTaskIDs, seqItems...)

	sequenceJson, err := json.Marshal(updatedSeqItems)
	if err != nil {
		return err
	}

	seqUpdateStmt, err := tx.Prepare(`
UPDATE task_sequence SET sequence = ? where id = 1;
`)
	if err != nil {
		return err
	}
	defer seqUpdateStmt.Close()

	_, err = seqUpdateStmt.Exec(sequenceJson)

	if err != nil {
		return err
	}

	err = tx.Commit()
	if err != nil {
		return err
	}
	return nil
}

func UpdateTaskSummaryInDB(db *sql.DB, id uint64, summary string) error {

	stmt, err := db.Prepare(`
UPDATE task
SET summary = ?
WHERE id = ?
`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(summary, id)

	if err != nil {
		return err
	}
	return nil
}

func ChangeTaskStatusInDB(db *sql.DB, id uint64, active bool, updatedAt time.Time) error {

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

func FetchActiveTasksFromDB(db *sql.DB, limit int) ([]types.Task, error) {

	var tasks []types.Task

	rows, err := db.Query(`
SELECT t.id, t.summary, t.created_at, t.updated_at
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
	return tasks, nil
}

func FetchInActiveTasksFromDB(db *sql.DB, limit int) ([]types.Task, error) {

	var tasks []types.Task

	rows, err := db.Query(`
SELECT id, summary, created_at, updated_at
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
	return tasks, nil
}

func DeleteTaskInDB(db *sql.DB, id uint64) error {

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
