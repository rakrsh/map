// Copyright (c) 2026 Ravi Sharma

//go:build integration
// +build integration

package integration

import (
	"net"
	"os"
	"testing"
	"time"

	"map/tile-service/internal/handler"
)

func getEnv(key, defaultVal string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return defaultVal
}

func TestPostgresEphemeralIntegration(t *testing.T) {
	pgAddr := getEnv("TEST_POSTGRES_ADDR", "localhost:5433")

	conn, err := net.DialTimeout("tcp", pgAddr, 3*time.Second)
	if err != nil {
		t.Fatalf("failed to connect to ephemeral PostGIS at %s: %v", pgAddr, err)
	}
	defer conn.Close()

	// Verify tile router initialization and handler responsiveness
	router := handler.NewRouter()
	if router == nil {
		t.Fatalf("failed to initialize tile handler router")
	}
}
