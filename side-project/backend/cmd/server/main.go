package main

import (
	"log"
	"net/http"
	"os"

	"exchange-travel-planner/backend/internal/httpapi"
	"exchange-travel-planner/backend/internal/store"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	st := store.New()
	srv := httpapi.NewServer(st)

	log.Printf("go backend running on :%s", port)
	if err := http.ListenAndServe(":"+port, srv.Routes()); err != nil {
		log.Fatal(err)
	}
}
