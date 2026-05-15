// Package docs provides a static Swagger UI for the Tasks API.
// No swag code generation is required — the spec is embedded directly.
//
// To regenerate from swag annotations run:
//
//	go install github.com/swaggo/swag/cmd/swag@latest
//	swag init -g cmd/server/main.go -o docs
package docs

import (
	_ "embed"
	"net/http"
)

//go:embed swagger.json
var swaggerJSON []byte

// SwaggerJSON returns the raw OpenAPI 2.0 specification.
func SwaggerJSON() []byte { return swaggerJSON }

// Handler returns an http.Handler that serves the Swagger UI.
// It mounts:
//
//	GET /swagger/         → Swagger UI (HTML)
//	GET /swagger/doc.json → OpenAPI JSON spec
func Handler() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("/swagger/doc.json", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		_, _ = w.Write(swaggerJSON)
	})

	mux.HandleFunc("/swagger/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		_, _ = w.Write([]byte(swaggerUIHTML))
	})

	return mux
}

const swaggerUIHTML = `<!DOCTYPE html>
<html>
<head>
  <meta charset="utf-8">
  <title>Tasks API — Swagger UI</title>
  <link rel="stylesheet" href="https://unpkg.com/swagger-ui-dist@5/swagger-ui.css">
</head>
<body>
  <div id="swagger-ui"></div>
  <script src="https://unpkg.com/swagger-ui-dist@5/swagger-ui-bundle.js"></script>
  <script>
    SwaggerUIBundle({
      url: '/swagger/doc.json',
      dom_id: '#swagger-ui',
      presets: [SwaggerUIBundle.presets.apis, SwaggerUIBundle.SwaggerUIStandalonePreset],
      layout: 'BaseLayout',
      deepLinking: true,
    });
  </script>
</body>
</html>`
