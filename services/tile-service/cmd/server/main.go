// Copyright (c) 2026 Ravi Sharma

package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"map/tile-service/internal/handler"
)

func main() {
	r := handler.NewRouter()

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
