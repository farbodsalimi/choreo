package statemachine

import (
	"fmt"
	"sync"

	_ "github.com/mattn/go-sqlite3"

	"github.com/farbodsalimi/choreo/internal/checkpoint"
	"github.com/farbodsalimi/choreo/internal/task"
)

type StateMachine struct {
	Tasks      map[string]task.Task
	State      map[string]any
	Completed  map[string]bool
	Checkpoint checkpoint.Checkpoint
	mu         sync.Mutex
	wg         sync.WaitGroup
}

func NewStateMachine(tasks []task.Task, checkpoint checkpoint.Checkpoint) (*StateMachine, error) {
	sm := &StateMachine{
		Tasks:      make(map[string]task.Task),
		State:      make(map[string]any),
		Completed:  make(map[string]bool),
		Checkpoint: checkpoint,
	}
	for _, task := range tasks {
		if _, exists := sm.Tasks[task.Name]; exists {
			return nil, fmt.Errorf("duplicate task name: %s", task.Name)
		}
		sm.Tasks[task.Name] = task
	}
	if err := sm.checkForCycles(); err != nil {
		return nil, err
	}
	sm.loadCheckpoint()
	return sm, nil
}

func (sm *StateMachine) checkForCycles() error {
	fmt.Println("Checking for cycles in tasks...")
	visited := make(map[string]bool)  //
	recStack := make(map[string]bool) //

	var dfs func(string) bool
	dfs = func(taskName string) bool {
		if recStack[taskName] {
			return true // Cycle detected
		}
		if visited[taskName] {
			return false
		}
		visited[taskName] = true
		recStack[taskName] = true
		for _, dep := range sm.Tasks[taskName].DependsOn {
			if dfs(dep) {
				return true
			}
		}
		recStack[taskName] = false
		return false
	}

	for taskName := range sm.Tasks {
		if !visited[taskName] {
			if dfs(taskName) {
				return fmt.Errorf("cycle detected in tasks on %s", taskName)
			}
		}
	}
	return nil
}

func (sm *StateMachine) loadCheckpoint() {
	fmt.Println("Loading checkpoint...")
	sm.Checkpoint.Load(&sm.State, &sm.Completed)
	fmt.Println("Checkpoint loaded.")
}

func (sm *StateMachine) saveCheckpoint() {
	fmt.Println("Saving checkpoint...")
	sm.Checkpoint.Save(sm.State, sm.Completed)
	fmt.Println("Checkpoint saved.")
}

func (sm *StateMachine) canRun(task task.Task) bool {
	for _, dep := range task.DependsOn {
		if !sm.Completed[dep] {
			return false
		}
	}
	return true
}

func (sm *StateMachine) runTask(name string, task task.Task) {
	defer sm.wg.Done()
	err := task.Func(sm.State)
	sm.mu.Lock()
	defer sm.mu.Unlock()
	if err == nil {
		sm.Completed[name] = true
		sm.saveCheckpoint()
	} else {
		fmt.Printf("Task %s failed: %v\n", name, err)
	}
}

func (sm *StateMachine) Run() error {
	for {
		var runnableTasks []string
		sm.mu.Lock()
		for name, task := range sm.Tasks {
			if !sm.Completed[name] && sm.canRun(task) {
				runnableTasks = append(runnableTasks, name)
			}
		}
		sm.mu.Unlock()

		if len(runnableTasks) == 0 {
			break
		}

		for _, name := range runnableTasks {
			sm.wg.Add(1)
			go sm.runTask(name, sm.Tasks[name])
		}

		sm.wg.Wait() // Wait for current batch of tasks to finish
	}
	return nil
}
