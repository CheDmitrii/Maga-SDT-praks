// Package main Notes API server.
//
// @title           Notes API
// @version         1.0
// @description     Учебный REST API для заметок (CRUD).
// @contact.name    Backend Course
// @contact.email   example@university.ru
// @BasePath        /api/v1
package main

import (
	"Prak_12/docs"
	_ "Prak_12/docs"
	httpx "Prak_12/internal/http"
	"Prak_12/internal/repo"
	httpSwagger "github.com/swaggo/http-swagger"
	"log"
	"net/http"
)

func main() {
	// Настроим SwaggerInfo
	docs.SwaggerInfo.Host = "localhost:8085"
	docs.SwaggerInfo.Schemes = []string{"http"}
	docs.SwaggerInfo.BasePath = "/api/v1"

	mem := repo.NewNoteRepoMem()

	r := httpx.NewRouter(mem)

	r.Get("/docs/*", httpSwagger.WrapHandler)

	log.Println("Server started at :8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
