package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"prak3/internal/api"
	"prak3/internal/storage"
	"strings"
	"testing"
)

func TestCreateTask_Valid(t *testing.T) {
	store := storage.NewMemoryStore()
	h := api.NewHandlers(store)

	body := bytes.NewBufferString(`{"title":"Test task"}`)
	req := httptest.NewRequest(http.MethodPost, "/tasks", body)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	h.CreateTask(w, req)

	if w.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d", w.Code)
	}

	var task map[string]any
	if err := json.Unmarshal(w.Body.Bytes(), &task); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}

	if task["title"] != "Test task" {
		t.Errorf("expected title 'Test task', got %v", task["title"])
	}
}

func TestCreateTask_TooShort(t *testing.T) {
	store := storage.NewMemoryStore()
	h := api.NewHandlers(store)

	body := bytes.NewBufferString(`{"title":"` + strings.Repeat("A", 145) + `"}`)
	req := httptest.NewRequest(http.MethodPost, "/tasks", body)
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	h.CreateTask(w, req)

	if w.Code != http.StatusUnprocessableEntity {
		t.Fatalf("expected 422, got %d", w.Code)
	}
}

func TestListTasks_Filter(t *testing.T) {
	store := storage.NewMemoryStore()
	store.Create("Buy milk")
	store.Create("Write code")

	h := api.NewHandlers(store)
	req := httptest.NewRequest(http.MethodGet, "/tasks?q=milk", nil)
	w := httptest.NewRecorder()

	h.ListTasks(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}

	var tasks []map[string]any
	if err := json.Unmarshal(w.Body.Bytes(), &tasks); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}

	if len(tasks) != 1 || tasks[0]["title"] != "Buy milk" {
		t.Errorf("filter failed, got %v", tasks)
	}
}
