package main

import (
	"fmt"
	"os"
	"strings"
	"time"
)

type Task struct {
	ID          int
	Description string
	Status      string
	CreatedAt   time.Time
	UpdatedAt   time.Time
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

func main() {
	if len(os.Args) < 2 {
		printUsage()
		return
	}

	tasks := []Task{}

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
	case "--help", "-h":
		printUsage()
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
