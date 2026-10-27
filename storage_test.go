package main

import (
	"os"
	"path/filepath"
	"testing"
)

func TestMissingFileIsCreatedAndCanBeRead(t *testing.T) {
	oldDir, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}

	tempDir := t.TempDir()
	if err := os.Chdir(tempDir); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		_ = os.Chdir(oldDir)
	})

	tasks, err := loadTasks()
	if err != nil {
		t.Fatal(err)
	}
	if len(tasks) != 0 {
		t.Fatalf("expected no tasks, got %d", len(tasks))
	}
	if _, err := os.Stat(filepath.Join(tempDir, dataFile)); err != nil {
		t.Fatal(err)
	}

	id, err := addTask(&tasks, "persist me")
	if err != nil {
		t.Fatal(err)
	}
	if id != 1 {
		t.Fatalf("ID = %d; want 1", id)
	}
	if err := saveTasks(tasks); err != nil {
		t.Fatal(err)
	}

	loaded, err := loadTasks()
	if err != nil {
		t.Fatal(err)
	}
	if len(loaded) != 1 || loaded[0].Description != "persist me" {
		t.Fatalf("unexpected tasks: %#v", loaded)
	}
}
