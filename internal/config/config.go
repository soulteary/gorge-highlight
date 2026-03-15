package config

import (
	"encoding/json"
	"os"
	"strconv"
)

type Config struct {
	ListenAddr   string `json:"listenAddr"`
	ServiceToken string `json:"serviceToken"`
	MaxBytes     int    `json:"maxBytes"`
	TimeoutSec   int    `json:"timeoutSec"`
}

func LoadFromEnv() *Config {
	return &Config{
		ListenAddr:   envStr("LISTEN_ADDR", ":8140"),
		ServiceToken: envStr("SERVICE_TOKEN", ""),
		MaxBytes:     envInt("MAX_BYTES", 1048576),
		TimeoutSec:   envInt("TIMEOUT_SEC", 15),
	}
}

func LoadFromFile(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	cfg := &Config{
		ListenAddr:   ":8140",
		ServiceToken: envStr("SERVICE_TOKEN", ""),
		MaxBytes:     1048576,
		TimeoutSec:   15,
	}
	if err := json.Unmarshal(data, cfg); err != nil {
		return nil, err
	}
	return cfg, nil
}

func envStr(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func envInt(key string, fallback int) int {
	if v := os.Getenv(key); v != "" {
		n, err := strconv.Atoi(v)
		if err == nil {
			return n
		}
	}
	return fallback
}
