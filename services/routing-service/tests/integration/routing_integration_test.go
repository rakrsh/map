// Copyright (c) 2026 Ravi Sharma

//go:build integration
// +build integration

package integration

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
	"testing"
	"time"

	"map/routing-service/internal/graph"
	"map/routing-service/internal/router"
)

func getEnv(key, defaultVal string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return defaultVal
}

func TestRedisEphemeralIntegration(t *testing.T) {
	redisAddr := getEnv("TEST_REDIS_ADDR", "localhost:6380")

	conn, err := net.DialTimeout("tcp", redisAddr, 3*time.Second)
	if err != nil {
		t.Fatalf("failed to connect to ephemeral Redis at %s: %v", redisAddr, err)
	}
	defer conn.Close()

	_ = conn.SetDeadline(time.Now().Add(3 * time.Second))

	// Send PING command
	fmt.Fprintf(conn, "PING\r\n")
	reader := bufio.NewReader(conn)
	response, err := reader.ReadString('\n')
	if err != nil {
		t.Fatalf("failed to read PING response from Redis: %v", err)
	}

	if !strings.Contains(response, "PONG") {
		t.Fatalf("expected PONG response, got %q", response)
	}

	// Verify route computation works alongside cache probe
	g := graph.Graph{
		Edges: []graph.Edge{
			{ID: 1, From: 10, To: 20, Weight: 12.5},
		},
	}
	path, err := router.FindShortestPath(g, 10, 20)
	if err != nil {
		t.Fatalf("router calculation failed: %v", err)
	}
	if len(path) != 2 || path[0] != 10 || path[1] != 20 {
		t.Errorf("unexpected path: %v", path)
	}
}
