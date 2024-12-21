package checkpoint

import (
	"encoding/json"
	"os"
)

var _ Checkpoint = &FileCheckpoint{}

type FileCheckpoint struct {
	FilePath string
}

func (fc *FileCheckpoint) Load(state *map[string]any, completed *map[string]bool) error {
	if _, err := os.Stat(fc.FilePath); err == nil {
		data, err := os.ReadFile(fc.FilePath)
		if err != nil {
			return err
		}
		type CheckpointFile struct {
			Completed map[string]bool `json:"completed"`
			State     map[string]any  `json:"state"`
		}
		var checkpointFile CheckpointFile
		json.Unmarshal(data, &checkpointFile)
		*state = checkpointFile.State
		*completed = checkpointFile.Completed
		return err
	}
	return nil
}

func (fc *FileCheckpoint) Save(state map[string]any, completed map[string]bool) error {
	data, err := json.Marshal(map[string]any{"state": state, "completed": completed})
	if err != nil {
		return err
	}
	return os.WriteFile(fc.FilePath, data, 0644)
}
