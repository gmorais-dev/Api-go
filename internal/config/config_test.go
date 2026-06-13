package config

import (
	"os"
	"testing"
	"time"
)

func TestNewAppConfigUsesEnvironmentOverrides(t *testing.T) {
	t.Setenv("SCRAPER_INPUT_FILE", "custom/input.txt")
	t.Setenv("SCRAPER_OUTPUT_FILE", "custom/output.csv")
	t.Setenv("SCRAPER_WORKER_COUNT", "7")
	t.Setenv("SCRAPER_JOB_TIMEOUT", "4s")
	t.Setenv("SCRAPER_GLOBAL_TIMEOUT", "45s")

	cfg, err := NewAppConfig()
	if err != nil {
		t.Fatal(err)
	}

	if cfg.InputFilePath != "custom/input.txt" {
		t.Fatalf("input inesperado: %s", cfg.InputFilePath)
	}

	if cfg.OutputFilePath != "custom/output.csv" {
		t.Fatalf("output inesperado: %s", cfg.OutputFilePath)
	}

	if cfg.WorkerCount != 7 {
		t.Fatalf("worker count inesperado: %d", cfg.WorkerCount)
	}

	if cfg.JobTimeout != 4*time.Second {
		t.Fatalf("job timeout inesperado: %s", cfg.JobTimeout)
	}

	if cfg.GlobalTimeout != 45*time.Second {
		t.Fatalf("global timeout inesperado: %s", cfg.GlobalTimeout)
	}
}

func TestNewAppConfigRejectsInvalidWorkerCount(t *testing.T) {
	t.Setenv("SCRAPER_WORKER_COUNT", "0")

	_, err := NewAppConfig()
	if err == nil {
		t.Fatal("esperava erro para worker count invalido")
	}
}

func TestGetEnvOrDefault(t *testing.T) {
	const key = "SCRAPER_SAMPLE_ENV"

	if err := os.Unsetenv(key); err != nil {
		t.Fatal(err)
	}

	if value := getEnvOrDefault(key, "fallback"); value != "fallback" {
		t.Fatalf("fallback inesperado: %s", value)
	}
}
