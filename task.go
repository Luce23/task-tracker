package main

import (
	"errors"
	"fmt"
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

func validateTask(task Task) error {
	if task.ID < 1 {
		return errors.New("ID должен быть положительным")
	}
	if strings.TrimSpace(task.Description) == "" {
		return errors.New("описание не может быть пустым")
	}
	if task.Status != "todo" && task.Status != "in-progress" && task.Status != "done" {
		return fmt.Errorf("недопустимый статус %q", task.Status)
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
