package task

type State map[string]any

func (s State) Set(key, value string) State {
	s[key] = value
	return s
}

type TaskFunc func(State) error

type Task struct {
	Name      string
	Func      TaskFunc
	DependsOn []string
}
