package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strconv"
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

func addTask(tasks *[]Task, description string) (int, error) {
	description = strings.TrimSpace(description)
	if description == "" {
		return 0, errors.New("описание не может быть пустым")
	}
	now := time.Now().UTC()
	task := Task{
		ID:          nextID(*tasks),
		Description: description,
		Status:      "todo",
		CreatedAt:   now,
		UpdatedAt:   now,
	}
	*tasks = append(*tasks, task)
	return task.ID, nil
}

func parseID(raw string) (int, error) {
	id, err := strconv.Atoi(raw)
	if err != nil || id < 1 {
		return 0, fmt.Errorf("ID должен быть положительным числом: %q", raw)
	}
	return id, nil
}

func findTask(tasks []Task, id int) (int, error) {
	for i := range tasks {
		if tasks[i].ID == id {
			return i, nil
		}
	}
	return -1, fmt.Errorf("Задача с ID %d не найдена", id)
}

func deleteTask(tasks *[]Task, id int) error {
	i, err := findTask(*tasks, id)
	if err != nil {
		return err
	}
	*tasks = append((*tasks)[:i], (*tasks)[i+1:]...)
	return nil

}

func updateTask(tasks []Task, id int, description string) error {
	description = strings.TrimSpace(description)
	if description == "" {
		return errors.New("Описание не может быть пустым.")
	}

	i, err := findTask(tasks, id)
	if err != nil {
		return err
	}

	tasks[i].Description = description
	tasks[i].UpdatedAt = time.Now().UTC()
	return nil
}

func markTask(tasks []Task, id int, status string) error {
	i, err := findTask(tasks, id)
	if err != nil {
		return err
	}

	tasks[i].Status = status
	tasks[i].UpdatedAt = time.Now().UTC()
	return nil
}

func filterTasks(tasks []Task, status string) ([]Task, error) {
	if status != "todo" && status != "in-progress" && status != "done" {
		return nil, fmt.Errorf("неизвестный статус %q", status)
	}

	result := make([]Task, 0)
	for _, task := range tasks {
		if task.Status == status {
			result = append(result, task)
		}
	}
	return result, nil
}

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
