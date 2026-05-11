package main

import (
	"log"
	"net/http"

	"Prak_4/internal/httpapi"
	"Prak_4/internal/student"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	repo := student.NewRepo()
	handler := httpapi.NewHandler(repo)

	mux := http.NewServeMux()
	mux.HandleFunc("/health", handler.Health)
	mux.HandleFunc("/students/", handler.GetStudentByID)
	// Expose Prometheus metrics — excluded from MetricsMiddleware to avoid self-counting
	mux.Handle("/metrics", promhttp.Handler())

	rootHandler := httpapi.MetricsMiddleware(mux)

	log.Println("server started on :8080")
	if err := http.ListenAndServe(":8080", rootHandler); err != nil {
		log.Fatal(err)
	}
}
