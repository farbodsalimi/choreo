package checkpoint

import (
	"fmt"
	"os"
	"testing"
)

func TestFileCheckpoint(t *testing.T) {
	fileCheckpoint := &FileCheckpoint{FilePath: "test_checkpoint.json"}
	defer os.Remove("test_checkpoint.json")

	state := map[string]any{"key": "value"}
	completed := map[string]bool{"Task1": true}

	// Save checkpoint
	if err := fileCheckpoint.Save(state, completed); err != nil {
		t.Fatalf("Failed to save checkpoint: %v", err)
	}

	// Load checkpoint
	loadedState := make(map[string]any)
	loadedCompleted := make(map[string]bool)
	if err := fileCheckpoint.Load(&loadedState, &loadedCompleted); err != nil {
		t.Fatalf("Failed to load checkpoint: %v", err)
	}

	fmt.Println(loadedState, loadedCompleted)

	// Verify loaded values
	if loadedState["key"] != "value" || !loadedCompleted["Task1"] {
		t.Errorf("Loaded checkpoint data mismatch")
	}
}
