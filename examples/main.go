package main

import (
	"fmt"

	"github.com/farbodsalimi/choreo/internal/checkpoint"
	"github.com/farbodsalimi/choreo/internal/statemachine"
	"github.com/farbodsalimi/choreo/internal/task"
)

func main() {
	tasks := []task.Task{
		{
			Name: "Counter1",
			Func: func(state task.State) error {
				state["counter"] = 100
				fmt.Println("Task1 completed")
				return nil
			},
			DependsOn: []string{},
		},
		{
			Name: "Counter2",
			Func: func(state task.State) error {

				c, _ := state["counter"].(int)
				state["counter"] = c * 2

				fmt.Println("Task2 completed")
				return nil
			},
			DependsOn: []string{"Counter1"},
		},
		{
			Name: "Counter3",
			Func: func(state task.State) error {
				c, _ := state["counter"].(int)
				state["counter"] = c / 10

				fmt.Println("Task3 completed")
				return nil
			},
			DependsOn: []string{"Counter2"},
		},
	}

	// checkpoint := &checkpoint.FileCheckpoint{FilePath: "checkpoint.json"}
	checkpoint, err := checkpoint.NewSQLiteCheckpoint("checkpoint.db")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	sm, err := statemachine.NewStateMachine(tasks, checkpoint)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	if err := sm.Run(); err != nil {
		fmt.Println("Error:", err)
	} else {
		fmt.Println("All tasks completed")
	}
}
