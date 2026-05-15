package http

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync/atomic"

	"Prak_13-producer/internal/amqp"
	"Prak_13-producer/internal/events"
)

type createTaskRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
}

type createTaskResponse struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
}

var seq atomic.Int64

type Handler struct {
	publisher *amqp.Publisher
}

func NewHandler(publisher *amqp.Publisher) *Handler {
	return &Handler{publisher: publisher}
}

func (h *Handler) CreateTask(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req createTaskRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "bad json", http.StatusBadRequest)
		return
	}

	if req.Title == "" {
		http.Error(w, "title is required", http.StatusBadRequest)
		return
	}

	n := seq.Add(1)
	taskID := fmt.Sprintf("t_%03d", n)

	requestID := r.Header.Get("X-Request-ID")

	// Publish event — best effort
	ev := events.NewTaskCreated(taskID, requestID)
	if err := h.publisher.Publish(context.Background(), ev); err != nil {
		log.Printf("publish failed (best effort): %v", err)
	}

	resp := createTaskResponse{
		ID:          taskID,
		Title:       req.Title,
		Description: req.Description,
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(resp)
}

func (h *Handler) Health(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	_ = json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}
