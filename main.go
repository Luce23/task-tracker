package main

import (
	"fmt"
	"os"
	"strings"
)

func printTasks(tasks []Task) {
	if len(tasks) == 0 {
		fmt.Println("Задач нет.")
		return
	}
	for _, task := range tasks {
		fmt.Printf("%d. [%s] %s\n", task.ID, task.Status, task.Description)
	}
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

		id, err := addTask(&tasks, strings.Join(os.Args[2:], " "))
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		if err := saveTasks(tasks); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		fmt.Printf("Задача добавлена (ID: %d)\n", id)

	case "list":
		if len(os.Args) == 2 {
			printTasks(tasks)
			return
		}
		if len(os.Args) != 3 {
			fmt.Fprintln(os.Stderr, "Использование: list [todo|in-progress|done]")
			os.Exit(1)
		}

		filtered, err := filterTasks(tasks, os.Args[2])
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		printTasks(filtered)
	case "get":
		if len(os.Args) != 3 {
			fmt.Fprintln(os.Stderr, "Использование: get <id>")
			os.Exit(1)
		}

		id, err := parseID(os.Args[2])
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		index, err := findTask(tasks, id)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		printTasks([]Task{tasks[index]})

	case "update":
		if len(os.Args) < 4 {
			fmt.Fprintln(os.Stderr, "Использование: update <id> <описание>")
			os.Exit(1)
		}
		id, err := parseID(os.Args[2])
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		if err := updateTask(tasks, id, strings.Join(os.Args[3:], " ")); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		if err := saveTasks(tasks); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		fmt.Printf("Задача %d обновлена\n", id)
	case "delete":
		if len(os.Args) != 3 {
			fmt.Fprintln(os.Stderr, "Использование: delete <id>")
			os.Exit(1)
		}
		id, err := parseID(os.Args[2])
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		if err := deleteTask(&tasks, id); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		if err := saveTasks(tasks); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		fmt.Printf("Задача %d удалена\n", id)

	case "mark-in-progress", "mark-done":
		if len(os.Args) != 3 {
			fmt.Fprintf(os.Stderr, "Использование: %s <id>\n", os.Args[1])
			os.Exit(1)
		}

		id, err := parseID(os.Args[2])
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}

		status := "done"
		if os.Args[1] == "mark-in-progress" {
			status = "in-progress"
		}

		if err := markTask(tasks, id, status); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		if err := saveTasks(tasks); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}

		fmt.Printf("Задача %d отмечена как %s\n", id, status)

	default:
		fmt.Fprintf(os.Stderr, "Неизвестная команда: %s\n", os.Args[1])
		printUsage()
		os.Exit(1)
	}
}

func printUsage() {
	fmt.Println(`Использование:
  task-cli add "описание"
  task-cli get <id>
  task-cli update <id> "новое описание"
  task-cli delete <id>
  task-cli mark-in-progress <id>
  task-cli mark-done <id>
  task-cli list [todo|in-progress|done]`)
}
