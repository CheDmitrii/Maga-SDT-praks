package events

import "time"

// TaskEvent is the message published to RabbitMQ after task creation.
type TaskEvent struct {
	Event     string `json:"event"`
	TaskID    string `json:"task_id"`
	TS        string `json:"ts"`
	RequestID string `json:"request_id,omitempty"`
	Producer  string `json:"producer,omitempty"`
}

func NewTaskCreated(taskID, requestID string) TaskEvent {
	return TaskEvent{
		Event:     "task.created",
		TaskID:    taskID,
		TS:        time.Now().UTC().Format(time.RFC3339),
		RequestID: requestID,
		Producer:  "producer-service",
	}
}
