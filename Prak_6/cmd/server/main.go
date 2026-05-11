package main

import (
	"log"
	"net/http"

	"Prak_6/internal/httpapi"
	"Prak_6/internal/store"
)

func main() {
	st := store.New()

	handler, err := httpapi.NewHandler(st)
	if err != nil {
		log.Fatal(err)
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/login", handler.Login)
	mux.HandleFunc("/logout", handler.Logout)
	mux.HandleFunc("/profile", handler.Profile)
	mux.HandleFunc("/hello", handler.Hello)

	log.Println("server started on http://localhost:8080")
	log.Println("open http://localhost:8080/login")

	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatal(err)
	}
}
