// Copyright (c) 2026 Ravi Sharma

package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
)

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/health", func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
	}).Methods(http.MethodGet)

	r.HandleFunc("/tiles/{z}/{x}/{y}", func(w http.ResponseWriter, req *http.Request) {
		vars := mux.Vars(req)
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]string{
			"message": "tile request",
			"z":       vars["z"],
			"x":       vars["x"],
			"y":       vars["y"],
		})
	}).Methods(http.MethodGet)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8082"
	}

	addr := fmt.Sprintf(":%s", port)
	log.Printf("tile-service listening on %s", addr)

	// Local Compose traffic is plain HTTP; TLS terminates at the deployment edge.
	// nosemgrep: go.lang.security.audit.net.use-tls.use-tls
	if err := http.ListenAndServe(addr, r); err != nil {
		log.Printf("server exited: %v", err)
		os.Exit(1)
	}
}
