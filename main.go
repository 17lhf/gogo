package main

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"gogo/internal/cache"
	"gogo/internal/config"
	"gogo/internal/db"
)

func main() {
	cfg := config.Load()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// PostgreSQL
	pool, err := db.New(ctx, cfg.Postgres.DSN())
	if err != nil {
		log.Fatalf("postgres init: %v", err)
	}
	defer pool.Close()

	if err := db.Ping(ctx, pool); err != nil {
		log.Fatalf("postgres unreachable: %v", err)
	}
	log.Println("postgres: connected")

	// Redis
	rdb := cache.New(cfg.Redis.Addr, cfg.Redis.Password)
	defer rdb.Close()

	if err := cache.Ping(context.Background(), rdb); err != nil {
		log.Fatalf("redis unreachable: %v", err)
	}
	log.Println("redis: connected")

	// Router
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"ok"}`))
	})

	log.Printf("HTTP listening on :%s", cfg.Server.Port)
	if err := http.ListenAndServe(":"+cfg.Server.Port, r); err != nil {
		log.Fatalf("server: %v", err)
	}
}
