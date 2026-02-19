package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"exchange-travel-planner/backend/internal/db"
	"exchange-travel-planner/backend/internal/domain"
	"exchange-travel-planner/backend/internal/httpapi"
	"exchange-travel-planner/backend/internal/store"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	var ds domain.DataStore

	if dsn := os.Getenv("DATABASE_URL"); dsn != "" {
		pool, err := db.Connect(context.Background(), dsn)
		if err != nil {
			log.Fatalf("postgres connect: %v", err)
		}
		ds = db.NewPgStore(pool)
		log.Println("using PostgreSQL store")
	} else {
		ds = store.New()
		log.Println("using in-memory store (set DATABASE_URL for Postgres)")
	}
	defer ds.Close()

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: httpapi.NewServer(ds).Routes(),
	}

	go func() {
		log.Printf("go backend running on :%s", port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal(err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("shutting down...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_ = srv.Shutdown(ctx)
}
