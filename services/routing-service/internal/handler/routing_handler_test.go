// Copyright (c) 2026 Ravi Sharma

package handler

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestHealthCheck(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	rec := httptest.NewRecorder()

	HealthCheck(rec, req)

	res := rec.Result()
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		t.Errorf("expected status 200, got %d", res.StatusCode)
	}

	contentType := res.Header.Get("Content-Type")
	if !strings.Contains(contentType, "application/json") {
		t.Errorf("expected application/json content type, got %s", contentType)
	}

	var body map[string]string
	if err := json.NewDecoder(res.Body).Decode(&body); err != nil {
		t.Fatalf("failed to decode response body: %v", err)
	}

	if body["status"] != "ok" {
		t.Errorf("expected status 'ok', got %q", body["status"])
	}
}

func TestGetRoute(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/routes", strings.NewReader(`{"start": [37.7749, -122.4194], "end": [37.7833, -122.4167]}`))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	GetRoute(rec, req)

	res := rec.Result()
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		t.Errorf("expected status 200, got %d", res.StatusCode)
	}

	var body map[string]string
	if err := json.NewDecoder(res.Body).Decode(&body); err != nil {
		t.Fatalf("failed to decode response body: %v", err)
	}

	if body["status"] != "pending_implementation" {
		t.Errorf("expected status 'pending_implementation', got %q", body["status"])
	}
	if body["message"] != "route request received" {
		t.Errorf("expected message 'route request received', got %q", body["message"])
	}
}
