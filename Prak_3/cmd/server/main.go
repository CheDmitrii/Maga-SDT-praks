package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"prak3/internal/api"
	"prak3/internal/storage"
	"time"
)

func main() {
	store := storage.NewMemoryStore()
	h := api.NewHandlers(store)

	mux := http.NewServeMux()
	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		api.JSON(w, http.StatusOK, map[string]string{"status": "ok"})
	})

	mux.HandleFunc("GET /tasks", h.ListTasks)
	mux.HandleFunc("POST /tasks", h.CreateTask)
	mux.HandleFunc("PATCH /tasks/", h.UpdateTask)
	mux.HandleFunc("DELETE /tasks/", h.DeleteTask)
	mux.HandleFunc("GET /tasks/", h.GetTask)

	handler := api.WithCORS(api.Logging(mux))
	addr := getAddr()
	server := &http.Server{
		Addr:    addr,
		Handler: handler,
	}

	log.Println("listening on", addr)
	if err := http.ListenAndServe(addr, handler); err != nil {
		log.Fatal(err)
	}

	// канал для прослушивания сигналов прерывания или завершения
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)

	go func() {
		log.Println("listening on", addr)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("ListenAndServe: %v", err)
		}
	}()

	<-stop // ожидание сигнала

	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server Shutdown Failed:%+v", err)
	}

	log.Println("Server exited gracefully")

}

func getAddr() string {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	return ":" + port
}
