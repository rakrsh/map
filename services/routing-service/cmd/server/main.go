// Copyright (c) 2026 Ravi Sharma

package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"

	"map/routing-service/internal/config"
	"map/routing-service/internal/handler"
)

func main() {
	cfg := config.Load()

	r := mux.NewRouter()
	r.HandleFunc("/health", handler.HealthCheck).Methods(http.MethodGet)
	r.HandleFunc("/routes", handler.GetRoute).Methods(http.MethodPost)

	addr := fmt.Sprintf(":%s", cfg.Port)
	log.Printf("routing-service listening on %s", addr)

	if err := http.ListenAndServe(addr, r); err != nil {
		log.Printf("server exited: %v", err)
		os.Exit(1)
	}
}
