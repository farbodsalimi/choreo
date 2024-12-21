package checkpoint

import (
	"os"
	"testing"
)

func TestSQLiteCheckpoint(t *testing.T) {
	sqliteCheckpoint, err := NewSQLiteCheckpoint("test_checkpoint.db")
	defer os.Remove("test_checkpoint.db")
	if err != nil {
		t.Fatalf("Failed to create SQLite checkpoint: %v", err)
	}

	state := map[string]any{"key": "value"}
	completed := map[string]bool{"Task1": true}

	// Save checkpoint
	if err := sqliteCheckpoint.Save(state, completed); err != nil {
		t.Fatalf("Failed to save checkpoint: %v", err)
	}

	// Load checkpoint
	loadedState := make(map[string]any)
	loadedCompleted := make(map[string]bool)
	if err := sqliteCheckpoint.Load(&loadedState, &loadedCompleted); err != nil {
		t.Fatalf("Failed to load checkpoint: %v", err)
	}

	// Verify loaded values
	if loadedState["key"] != "value" || !loadedCompleted["Task1"] {
		t.Errorf("Loaded checkpoint data mismatch")
	}
}
