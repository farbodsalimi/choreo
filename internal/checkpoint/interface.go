package checkpoint

type Checkpoint interface {
	Load(state *map[string]any, completed *map[string]bool) error
	Save(state map[string]any, completed map[string]bool) error
}
