package main

import (
	"log"
	"net/http"

	"Prak_1_user-service/internal/user"
)

func main() {
	repo := user.NewRepo()
	handler := user.NewHandler(repo)

	mux := http.NewServeMux()
	mux.HandleFunc("/users", handler.GetAllUsers)
	mux.HandleFunc("/users/", handler.GetUserByID)

	addr := ":8081"
	log.Println("user-service started on", addr)
	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatal(err)
	}
}
