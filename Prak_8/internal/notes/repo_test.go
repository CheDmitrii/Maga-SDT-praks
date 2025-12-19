package notes

import (
	"Prak_8/internal/db"
	"context"
	"os"
	"testing"
)

func TestCreateAndGet(t *testing.T) {
	ctx := context.Background()
	uri := getenv("MONGO_URI", "mongodb://root:secret@localhost:27017/?authSource=admin")
	dbName := "prak_8_test"
	deps, err := db.ConnectMongo(ctx, uri, dbName)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		deps.Client.Database(dbName).Drop(ctx)
		deps.Client.Disconnect(ctx)
	})
	r, err := NewRepo(deps.Database)
	if err != nil {
		t.Fatal(err)
	}

	created, err := r.Create(ctx, "T1", "C1", nil)
	if err != nil {
		t.Fatal(err)
	}

	got, err := r.ByID(ctx, created.ID.Hex())
	if err != nil {
		t.Fatal(err)
	}
	if got.Title != "T1" {
		t.Fatalf("want T1 got %s", got.Title)
	}
}

func getenv(k, def string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return def
}
