package main

import (
	"context"
	"log"
	"net/http"
	"strings"

	"Prak_9/internal/cache"
	"Prak_9/internal/config"
	"Prak_9/internal/httpapi"
	"Prak_9/internal/service"
	"Prak_9/internal/task"
)

func main() {
	cfg := config.New()

	repo := task.NewRepo()
	redisClient := cache.NewRedisClient(cfg)

	if err := cache.Ping(context.Background(), redisClient); err != nil {
		log.Println("warning: redis unavailable at startup:", err)
		log.Println("service will work without cache (fallback to repo)")
	} else {
		log.Println("redis connected at", cfg.RedisAddr)
	}

	taskService := service.NewTaskService(repo, redisClient, cfg)
	handler := httpapi.NewHandler(taskService)

	mux := http.NewServeMux()
	mux.HandleFunc("/v1/tasks/", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			id := strings.TrimPrefix(r.URL.Path, "/v1/tasks/")
			if id == "" {
				handler.GetTasks(w, r)
			} else {
				handler.GetTaskByID(w, r)
			}
		case http.MethodPatch:
			handler.PatchTask(w, r)
		case http.MethodDelete:
			handler.DeleteTask(w, r)
		default:
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		}
	})

	log.Println("server started on :8082")
	if err := http.ListenAndServe(":8082", mux); err != nil {
		log.Fatal(err)
	}
}
