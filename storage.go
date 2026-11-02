package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
)

const dataFile = "tasks.json"

func loadTasks() ([]Task, error) {
	data, err := os.ReadFile(dataFile)
	if errors.Is(err, os.ErrNotExist) {
		tasks := []Task{}
		if err := saveTasks(tasks); err != nil {
			return nil, err
		}
		return tasks, nil
	}

	if err != nil {
		return nil, fmt.Errorf("read %s: %w", dataFile, err)
	}

	var tasks []Task
	if err := json.Unmarshal(data, &tasks); err != nil {
		return nil, fmt.Errorf("decode %s: %w", dataFile, err)
	}
	if tasks == nil {
		tasks = []Task{}
	}
	return tasks, nil
}

func saveTasks(tasks []Task) error {
	for _, task := range tasks {
		if err := validateTask(task); err != nil {
			return fmt.Errorf("задача %d некорректна: %w", task.ID, err)
		}
	}
	data, err := json.MarshalIndent(tasks, "", " ")
	if err != nil {
		return fmt.Errorf("encode tasks: %w", err)
	}
	if err := os.WriteFile(dataFile, append(data, '\n'), 0644); err != nil {
		return fmt.Errorf("write %s: %w", dataFile, err)
	}
	return nil
}
