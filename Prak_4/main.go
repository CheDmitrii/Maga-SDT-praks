package main

import (
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/go-chi/chi/v5"
	chimw "github.com/go-chi/chi/v5/middleware"

	"prak4/internal/task"
	myMW "prak4/pkg/middleware"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = ":8080"
	} else if !strings.HasPrefix(port, ":") {
		port = ":" + port
	}

	repo := task.NewRepo()
	h := task.NewHandler(repo)

	r := chi.NewRouter()
	r.Use(chimw.RequestID)
	r.Use(chimw.Recoverer)
	r.Use(myMW.Logger)
	r.Use(myMW.SimpleCORS)

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("OK"))
	})

	// v1 без фильтров и без пагинаци
	r.Route("/api/v1", func(api chi.Router) {
		api.Mount("/tasks", h.RoutesV1())
	})

	// v2 с фильтрами и пагинацией
	r.Route("/api/v2", func(api chi.Router) {
		api.Mount("/tasks", h.RoutesV2())
	})

	//addr := ":8080"
	log.Printf("listening on %s", port)
	log.Fatal(http.ListenAndServe(port, r))
}
