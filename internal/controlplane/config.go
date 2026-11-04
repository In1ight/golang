package controlplane

import (
	"fmt"
	"os"
	"time"
)

type Config struct {
	HTTPAddr          string
	LogLevel          string
	ReadTimeout       time.Duration
	WriteTimeout      time.Duration
	ReadHeaderTimeout time.Duration
	IdleTimeout       time.Duration
	ShutdownTimeout   time.Duration
}

func parseTimeout(name string, defaultValue string) (time.Duration, error) {
	timeout := os.Getenv(name)
	if timeout == "" {
		timeout = defaultValue
	}

	timeoutDuration, err := time.ParseDuration(timeout)

	if err != nil {
		return 0, fmt.Errorf("invalid %s: %w", name, err)
	}

	if timeoutDuration <= 0 {
		return 0, fmt.Errorf("%s must be positive", name)
	}

	return timeoutDuration, nil
}

func LoadConfig() (Config, error) {
	addr := os.Getenv("HTTP_ADDR")

	if addr == "" {
		addr = ":8080"
	}

	level := os.Getenv("LOG_LEVEL")

	if level == "" {
		level = "info"
	}

	if level != "debug" &&
		level != "info" &&
		level != "warn" &&
		level != "error" {

		return Config{}, fmt.Errorf("invalid LOG_LEVEL: %s", level)
	}

	rTimeout, err := parseTimeout("READ_TIMEOUT", "5s")
	if err != nil {
		return Config{}, err
	}

	wTimeout, err := parseTimeout("WRITE_TIMEOUT", "10s")
	if err != nil {
		return Config{}, err
	}

	rhTimeout, err := parseTimeout("READ_HEADER_TIMEOUT", "2s")
	if err != nil {
		return Config{}, err
	}

	shTimeout, err := parseTimeout("SHUTDOWN_TIMEOUT", "8s")
	if err != nil {
		return Config{}, err
	}

	idleTimeout, err := parseTimeout("IDLE_TIMEOUT", "60s")
	if err != nil {
		return Config{}, err
	}

	cfg := Config{
		HTTPAddr:          addr,
		LogLevel:          level,
		ReadTimeout:       rTimeout,
		WriteTimeout:      wTimeout,
		ReadHeaderTimeout: rhTimeout,
		IdleTimeout:       idleTimeout,
		ShutdownTimeout:   shTimeout,
	}
	return cfg, nil
}
