package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"

	"Prak_14/internal/config"
	httptransport "Prak_14/internal/http"
	"Prak_14/internal/storage/postgres"
	rediscache "Prak_14/internal/storage/redis"
)

func main() {
	_ = godotenv.Load()

	cfg := config.FromEnv()
	dbUri := os.Getenv("DATABASE_URL")
	if dbUri == "" {
		dbUri = "postgres://postgres:postgres@localhost:5432/notes?sslmode=disable"
	}

	pgxCfg, err := pgxpool.ParseConfig(dbUri)
	if err != nil {
		log.Fatal(err)
	}
	pgxCfg.MaxConns = 20
	pgxCfg.MinConns = 5
	pgxCfg.MaxConnLifetime = time.Hour
	pgxCfg.ConnConfig.StatementCacheCapacity = 256

	pool, err := pgxpool.NewWithConfig(context.Background(), pgxCfg)
	if err != nil {
		log.Fatal(err)
	}
	defer pool.Close()

	repo := postgres.NewRepo(pool)

	// Redis cache
	cache, err := rediscache.New(cfg.RedisAddr, cfg.RedisPassword, cfg.RedisDB, cfg.CacheTTL)
	if err != nil {
		log.Fatal(err)
	}
	defer cache.Close()

	srv := httptransport.NewServer(repo, cache)

	log.Printf("listening on %s", cfg.HTTPAddr)
	if err := http.ListenAndServe(cfg.HTTPAddr, srv.Router()); err != nil {
		log.Fatal(err)
	}
}
