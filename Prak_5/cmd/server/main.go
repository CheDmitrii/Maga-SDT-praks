package main

import (
	"crypto/tls"
	"database/sql"
	"log"
	"net/http"

	_ "github.com/lib/pq"

	"Prak_5/internal/config"
	"Prak_5/internal/httpapi"
	"Prak_5/internal/student"
)

func main() {
	cfg := config.New()

	db, err := sql.Open("postgres", cfg.DSN)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatal("cannot connect to database:", err)
	}

	repo := student.NewRepo(db)

	stmt, err := repo.PrepareGetByID()
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()

	handler := httpapi.NewHandler(repo, stmt)

	mux := http.NewServeMux()
	mux.HandleFunc("/health", handler.Health)
	mux.HandleFunc("/students", handler.GetStudentByID)

	cert, err := tls.LoadX509KeyPair(cfg.CertFile, cfg.KeyFile)
	if err != nil {
		log.Fatal("failed to load TLS cert/key:", err)
	}

	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{cert},
	}

	httpsServer := &http.Server{
		Addr:      cfg.Addr, // :8443
		Handler:   mux,
		TLSConfig: tlsConfig,
	}

	// HTTP-сервер на порту 8080 — только редирект на HTTPS
	httpServer := &http.Server{
		Addr: ":8080",
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			target := "https://localhost:8443" + r.URL.RequestURI()
			http.Redirect(w, r, target, http.StatusMovedPermanently)
		}),
	}

	// Запускаем HTTP-редирект в фоне
	go func() {
		log.Println("HTTP redirect server started on http://localhost:8080  →  https://localhost:8443")
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal("HTTP server error:", err)
		}
	}()

	log.Println("HTTPS server started on https://localhost:8443")
	if err := httpsServer.ListenAndServeTLS("", ""); err != nil {
		log.Fatal(err)
	}
}
