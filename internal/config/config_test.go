package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadFromEnvDefaults(t *testing.T) {
	cfg := LoadFromEnv()

	if cfg.ListenAddr != ":8140" {
		t.Errorf("expected :8140, got %s", cfg.ListenAddr)
	}
	if cfg.MaxBytes != 1048576 {
		t.Errorf("expected 1048576, got %d", cfg.MaxBytes)
	}
	if cfg.TimeoutSec != 15 {
		t.Errorf("expected 15, got %d", cfg.TimeoutSec)
	}
}

func TestLoadFromEnvOverride(t *testing.T) {
	t.Setenv("LISTEN_ADDR", ":9999")
	t.Setenv("SERVICE_TOKEN", "tok123")
	t.Setenv("MAX_BYTES", "2048")
	t.Setenv("TIMEOUT_SEC", "30")

	cfg := LoadFromEnv()

	if cfg.ListenAddr != ":9999" {
		t.Errorf("expected :9999, got %s", cfg.ListenAddr)
	}
	if cfg.ServiceToken != "tok123" {
		t.Errorf("expected tok123, got %s", cfg.ServiceToken)
	}
	if cfg.MaxBytes != 2048 {
		t.Errorf("expected 2048, got %d", cfg.MaxBytes)
	}
	if cfg.TimeoutSec != 30 {
		t.Errorf("expected 30, got %d", cfg.TimeoutSec)
	}
}

func TestLoadFromFile(t *testing.T) {
	content := `{"listenAddr":":7777","maxBytes":512,"timeoutSec":5}`
	dir := t.TempDir()
	path := filepath.Join(dir, "config.json")
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	cfg, err := LoadFromFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if cfg.ListenAddr != ":7777" {
		t.Errorf("expected :7777, got %s", cfg.ListenAddr)
	}
	if cfg.MaxBytes != 512 {
		t.Errorf("expected 512, got %d", cfg.MaxBytes)
	}
}

func TestLoadFromFileNotFound(t *testing.T) {
	_, err := LoadFromFile("/nonexistent/config.json")
	if err == nil {
		t.Error("expected error for missing file")
	}
}
