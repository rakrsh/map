// Copyright (c) 2026 Ravi Sharma

package config

import (
	"os"
	"testing"
)

func TestLoadDefaults(t *testing.T) {
	origPort := os.Getenv("PORT")
	defer func() {
		if origPort != "" {
			_ = os.Setenv("PORT", origPort)
		} else {
			_ = os.Unsetenv("PORT")
		}
	}()

	_ = os.Unsetenv("PORT")
	cfg := Load()
	if cfg.Port != "8081" {
		t.Errorf("expected default port 8081, got %q", cfg.Port)
	}
}

func TestLoadCustomPort(t *testing.T) {
	origPort := os.Getenv("PORT")
	defer func() {
		if origPort != "" {
			_ = os.Setenv("PORT", origPort)
		} else {
			_ = os.Unsetenv("PORT")
		}
	}()

	_ = os.Setenv("PORT", "9999")
	cfg := Load()
	if cfg.Port != "9999" {
		t.Errorf("expected custom port 9999, got %q", cfg.Port)
	}
}
