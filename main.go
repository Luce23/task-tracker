package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strings"
	"time"
)

type Task struct {
	ID          int       `json:"id"`
	Description string    `json:"description"`
	Status      string    `json:"status"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

func nextID(tasks []Task) int {
	maxID := 0
	for _, task := range tasks {
		if task.ID > maxID {
			maxID = task.ID
		}
	}
	return maxID + 1
}

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
	data, err := json.MarshalIndent(tasks, "", " ")
	if err != nil {
		return fmt.Errorf("encode tasks: %w", err)
	}
	if err := os.WriteFile(dataFile, append(data, '\n'), 0644); err != nil {
		return fmt.Errorf("write %s: %w", dataFile, err)
	}
	return nil
}

func main() {
	if len(os.Args) < 2 || os.Args[1] == "--help" || os.Args[1] == "-h" {
		printUsage()
		return
	}

	tasks, err := loadTasks()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	switch os.Args[1] {
	case "add":
		if len(os.Args) < 3 {
			fmt.Fprintln(os.Stderr, "Укажите описание")
			os.Exit(1)
		}
		now := time.Now().UTC()

		task := Task{
			ID:          nextID(tasks),
			Description: strings.Join(os.Args[2:], " "),
			Status:      "todo",
			CreatedAt:   now,
			UpdatedAt:   now,
		}
		tasks = append(tasks, task)

		fmt.Printf("Задача добавлена (ID: %d, статус: %s)\n", task.ID, task.Status)

	case "list":
		if len(tasks) == 0 {
			fmt.Println("Задач нет.")
			return
		}
		for _, task := range tasks {
			fmt.Printf("%d. [%s] %s\n", task.ID, task.Status, task.Description)
		}
	case "update", "delete", "mark-in-progress", "mark-done":
		fmt.Printf("Команда %q распознана\n", os.Args[1])
	default:
		fmt.Fprintf(os.Stderr, "Неизвестная команда: %s\n", os.Args[1])
		printUsage()
		os.Exit(1)
	}
}

func printUsage() {
	fmt.Println(`Использование:
  task-cli add "описание"
  task-cli update <id> "новое описание"
  task-cli delete <id>
  task-cli mark-in-progress <id>
  task-cli mark-done <id>
  task-cli list [todo|in-progress|done]`)
}
