package main

import (
	"testing"
	"time"
)

func TestAddTaskTrimsDescriptionAndAssignsID(t *testing.T) {
	tasks := []Task{}

	id, err := addTask(&tasks, "  learn Go  ")
	if err != nil {
		t.Fatal(err)
	}
	if id != 1 || tasks[0].Description != "learn Go" || tasks[0].Status != "todo" {
		t.Fatalf("unexpected task: %#v", tasks)
	}

	id, err = addTask(&tasks, "second task")
	if err != nil {
		t.Fatal(err)
	}
	if id != 2 {
		t.Fatalf("second ID = %d; want 2", id)
	}
}

func TestParseIDRejectsInvalidValues(t *testing.T) {
	for _, value := range []string{"0", "-2", "abc"} {
		if _, err := parseID(value); err == nil {
			t.Errorf("parseID(%q) should fail", value)
		}
	}
}

func TestUpdateTaskChangesDescriptionAndTimestamp(t *testing.T) {
	oldTime := time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)
	tasks := []Task{{ID: 1, Description: "old", Status: "todo", UpdatedAt: oldTime}}

	if err := updateTask(tasks, 1, "  new  "); err != nil {
		t.Fatal(err)
	}
	if tasks[0].Description != "new" {
		t.Fatalf("description = %q; want %q", tasks[0].Description, "new")
	}
	if !tasks[0].UpdatedAt.After(oldTime) {
		t.Fatal("UpdatedAt was not advanced")
	}
}

func TestDeleteTaskRemovesMatchingID(t *testing.T) {
	tasks := []Task{{ID: 1}, {ID: 2}, {ID: 3}}

	if err := deleteTask(&tasks, 2); err != nil {
		t.Fatal(err)
	}
	if len(tasks) != 2 || tasks[1].ID != 3 {
		t.Fatalf("unexpected tasks after delete: %#v", tasks)
	}
}

func TestMarkTaskChangesStatus(t *testing.T) {
	tasks := []Task{{ID: 1, Status: "todo"}}

	if err := markTask(tasks, 1, "done"); err != nil {
		t.Fatal(err)
	}
	if tasks[0].Status != "done" {
		t.Fatalf("status = %q; want done", tasks[0].Status)
	}
}

func TestFilterTasks(t *testing.T) {
	tasks := []Task{{ID: 1, Status: "todo"}, {ID: 2, Status: "done"}}

	got, err := filterTasks(tasks, "done")
	if err != nil {
		t.Fatal(err)
	}
	if len(got) != 1 || got[0].ID != 2 {
		t.Fatalf("unexpected result: %#v", got)
	}
	if _, err := filterTasks(tasks, "unknown"); err == nil {
		t.Fatal("filterTasks should reject an unknown status")
	}
}
