package integration

import (
	"Prak_16/internal/db"
	"Prak_16/internal/httpapi"
	"Prak_16/internal/repo"
	"Prak_16/internal/service"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	_ "github.com/lib/pq"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
)

func withPostgres(t *testing.T) (dsn string, term func()) {
	t.Helper()
	ctx := context.Background()
	pg, err := postgres.RunContainer(ctx,
		postgres.WithDatabase("notes_test"),
		postgres.WithUsername("test"),
		postgres.WithPassword("test"),
	)
	if err != nil {
		t.Skipf("cannot start testcontainer: %v", err)
	}
	// host, _ := pg.Host(ctx)
	// port, _ := pg.MappedPort(ctx, "5432")
	// dsn = fmt.Sprintf("postgres://test:test@%s:%s/notes_test?sslmode=disable", host, port.Port())
	dsn = "postgres://test:test@localhost:54321/notes_test?sslmode=disable"
	return dsn, func() { _ = pg.Terminate(ctx) }
}

func Test_CreateAndGet_withTC(t *testing.T) {
	dsn, stop := withPostgres(t)
	defer stop()

	dbx, err := sql.Open("postgres", dsn)
	if err != nil {
		t.Fatal(err)
	}
	db.MustApplyMigrations(dbx)

	r := gin.Default()
	svc := service.Service{Notes: repo.NoteRepo{DB: dbx}}
	httpapi.Router{Svc: &svc}.Register(r)
	srv := httptest.NewServer(r)
	defer srv.Close()

	// Create
	resp, err := http.Post(srv.URL+"/notes", "application/json",
		strings.NewReader(`{"title":"Hello","content":"World"}`))
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != http.StatusCreated {
		t.Fatalf("status=%d want=201", resp.StatusCode)
	}
	var created map[string]any
	b, _ := io.ReadAll(resp.Body)
	_ = resp.Body.Close()
	_ = json.Unmarshal(b, &created)
	id := int64(created["id"].(float64))

	// Get
	resp2, err := http.Get(fmt.Sprintf("%s/notes/%d", srv.URL, id))
	if err != nil {
		t.Fatal(err)
	}
	if resp2.StatusCode != http.StatusOK {
		t.Fatalf("status=%d want=200", resp2.StatusCode)
	}
}

//func Test_CreateAndGet_withTC(t *testing.T) {
//
//}
