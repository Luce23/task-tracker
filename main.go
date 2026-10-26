package main

import (
	"fmt"
	"os"
	"strings"
)

type Task struct {
	ID          int
	Description string
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

		task := Task{
			ID:          len(tasks) + 1,
			Description: strings.Join(os.Args[2:], " "),
		}
		tasks = append(tasks, task)

		fmt.Printf("Задача добавлена (ID: %d)\n", task.ID)

	case "list":
		fmt.Printf("Задач в памяти: %d\n", len(tasks))
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
