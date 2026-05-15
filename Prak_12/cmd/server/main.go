// @title           Tasks API
// @version         1.0
// @description     REST + GraphQL сервис управления задачами. Практическое занятие №12.
// @BasePath        /
package main

import (
	"Prak_12/graph"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"log"
	"net/http"

	"Prak_12/docs"
	"Prak_12/internal/rest"
	"Prak_12/internal/store"
)

func main() {
	s := store.New()

	restHandler := rest.NewHandler(s)
	swaggerHandler := docs.Handler()

	mux := http.NewServeMux()

	// REST API — /v1/tasks и /v1/tasks/{id}
	mux.Handle("/v1/tasks", restHandler)
	mux.Handle("/v1/tasks/", restHandler)

	// GraphQL сервер
	gqlSrv := handler.NewDefaultServer(
		graph.NewExecutableSchema(graph.Config{
			Resolvers: &graph.Resolver{Store: s},
		}),
	)
	// GraphQL API — /v1/graphql
	// GET  → Playground
	// POST → GraphQL endpoint
	mux.HandleFunc("/v1/graphql", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			playground.Handler("GraphQL", "/v1/graphql")(w, r)
			return
		}
		gqlSrv.ServeHTTP(w, r)
	})

	// Swagger UI — http://localhost:8080/swagger/
	// Spec JSON  — http://localhost:8080/swagger/doc.json
	mux.Handle("/swagger/", swaggerHandler)

	// Корень — краткая справка по маршрутам
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		_, _ = w.Write([]byte(`Tasks API — ПЗ 12
=================
REST:
  GET    /v1/tasks          — список задач
  POST   /v1/tasks          — создать задачу
  GET    /v1/tasks/{id}     — получить задачу
  PATCH  /v1/tasks/{id}     — обновить задачу
  DELETE /v1/tasks/{id}     — удалить задачу

GraphQL:
  GET  /v1/graphql          — Playground
  POST /v1/graphql          — GraphQL endpoint

Документация:
  GET  /swagger/            — Swagger UI
  GET  /swagger/doc.json    — OpenAPI 2.0 JSON
`))
	})

	log.Println("server started on http://localhost:8080")
	log.Println("  REST:    http://localhost:8080/v1/tasks")
	log.Println("  GraphQL: http://localhost:8080/v1/graphql")
	log.Println("  Swagger: http://localhost:8080/swagger/")
	log.Fatal(http.ListenAndServe(":8080", mux))
}
