package checkpoint

import (
	"database/sql"
	"encoding/json"
)

var _ Checkpoint = &SQLiteCheckpoint{}

type SQLiteCheckpoint struct {
	DB *sql.DB
}

func NewSQLiteCheckpoint(dbPath string) (*SQLiteCheckpoint, error) {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, err
	}
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS checkpoint (
			id INTEGER PRIMARY KEY,
			state TEXT,
			completed TEXT
		);
	`)
	if err != nil {
		return nil, err
	}
	return &SQLiteCheckpoint{DB: db}, nil
}

func (sc *SQLiteCheckpoint) Load(state *map[string]any, completed *map[string]bool) error {
	row := sc.DB.QueryRow("SELECT state, completed FROM checkpoint WHERE id = 1")
	var stateJSON, completedJSON string
	if err := row.Scan(&stateJSON, &completedJSON); err != nil {
		if err == sql.ErrNoRows {
			return nil // No data yet, so nothing to load
		}
		return err
	}
	if err := json.Unmarshal([]byte(stateJSON), state); err != nil {
		return err
	}
	return json.Unmarshal([]byte(completedJSON), completed)
}

func (sc *SQLiteCheckpoint) Save(state map[string]any, completed map[string]bool) error {
	stateJSON, err := json.Marshal(state)
	if err != nil {
		return err
	}
	completedJSON, err := json.Marshal(completed)
	if err != nil {
		return err
	}
	_, err = sc.DB.Exec(`
		INSERT INTO checkpoint (id, state, completed)
		VALUES (1, ?, ?)
		ON CONFLICT(id) DO UPDATE SET state=excluded.state, completed=excluded.completed;
	`, string(stateJSON), string(completedJSON))
	return err
}
