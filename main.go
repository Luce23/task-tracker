package main

import (
	"fmt"
	"os"
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

	switch os.Args[1] {
	case "list":
		tasks := []Task{}
		fmt.Printf("Задач в памяти: %d\n", len(tasks))
	case "add", "update", "delete", "mark-in-progress", "mark-done":
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
