package http

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	amqpclient "Prak_14-producer/internal/amqp"
	"Prak_14-producer/internal/jobs"
	amqplib "github.com/rabbitmq/amqp091-go"

	"crypto/rand"
	"encoding/hex"
)

type Handler struct {
	ch *amqplib.Channel
}

func NewHandler(ch *amqplib.Channel) *Handler {
	return &Handler{ch: ch}
}

func randomID() string {
	b := make([]byte, 8)
	_, _ = rand.Read(b)
	return hex.EncodeToString(b)
}

func (h *Handler) EnqueueJob(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		TaskID string `json:"task_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.TaskID == "" {
		http.Error(w, "task_id is required", http.StatusBadRequest)
		return
	}

	job := jobs.TaskJob{
		Job:       "process_task",
		TaskID:    req.TaskID,
		Attempt:   1,
		MessageID: fmt.Sprintf("msg_%s", randomID()),
	}

	if err := amqpclient.PublishJob(h.ch, "task_jobs", job); err != nil {
		log.Printf("enqueue error: %v", err)
		http.Error(w, "failed to enqueue job", http.StatusInternalServerError)
		return
	}

	log.Printf("job enqueued: %+v", job)

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusAccepted)
	_ = json.NewEncoder(w).Encode(map[string]string{
		"status":  "accepted",
		"task_id": req.TaskID,
	})
}

func (h *Handler) Health(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	_ = json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}
