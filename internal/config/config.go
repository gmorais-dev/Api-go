package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

type AppConfig struct {
	InputFilePath  string
	OutputFilePath string
	WorkerCount    int
	JobTimeout     time.Duration
	GlobalTimeout  time.Duration
}

const (
	defaultInputFilePath  = "input/urls.txt"
	defaultOutputFilePath = "output/results.csv"
	defaultWorkerCount    = 3
	defaultJobTimeout     = 2 * time.Second
	defaultGlobalTimeout  = 30 * time.Second
)

func NewAppConfig() (AppConfig, error) {
	cfg := AppConfig{
		InputFilePath:  getEnvOrDefault("SCRAPER_INPUT_FILE", defaultInputFilePath),
		OutputFilePath: getEnvOrDefault("SCRAPER_OUTPUT_FILE", defaultOutputFilePath),
		WorkerCount:    defaultWorkerCount,
		JobTimeout:     defaultJobTimeout,
		GlobalTimeout:  defaultGlobalTimeout,
	}

	workerCount, err := getEnvInt("SCRAPER_WORKER_COUNT", defaultWorkerCount)
	if err != nil {
		return AppConfig{}, err
	}
	cfg.WorkerCount = workerCount

	jobTimeout, err := getEnvDuration("SCRAPER_JOB_TIMEOUT", defaultJobTimeout)
	if err != nil {
		return AppConfig{}, err
	}
	cfg.JobTimeout = jobTimeout

	globalTimeout, err := getEnvDuration("SCRAPER_GLOBAL_TIMEOUT", defaultGlobalTimeout)
	if err != nil {
		return AppConfig{}, err
	}
	cfg.GlobalTimeout = globalTimeout

	if err := cfg.Validate(); err != nil {
		return AppConfig{}, err
	}

	return cfg, nil
}

func (cfg AppConfig) Validate() error {
	if cfg.InputFilePath == "" {
		return fmt.Errorf("o caminho do arquivo de entrada nao pode ser vazio")
	}

	if cfg.OutputFilePath == "" {
		return fmt.Errorf("o caminho do arquivo de saida nao pode ser vazio")
	}

	if cfg.WorkerCount <= 0 {
		return fmt.Errorf("worker count invalido: %d", cfg.WorkerCount)
	}

	if cfg.JobTimeout <= 0 {
		return fmt.Errorf("job timeout invalido: %s", cfg.JobTimeout)
	}

	if cfg.GlobalTimeout <= 0 {
		return fmt.Errorf("global timeout invalido: %s", cfg.GlobalTimeout)
	}

	return nil
}

func getEnvOrDefault(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		value = strings.TrimSpace(value)
		if value != "" {
			return value
		}
	}

	return fallback
}

func getEnvInt(key string, fallback int) (int, error) {
	value := strings.TrimSpace(os.Getenv(key))
	if value == "" {
		return fallback, nil
	}

	parsed, err := strconv.Atoi(value)
	if err != nil {
		return 0, fmt.Errorf("valor invalido para %s: %w", key, err)
	}

	return parsed, nil
}

func getEnvDuration(key string, fallback time.Duration) (time.Duration, error) {
	value := strings.TrimSpace(os.Getenv(key))
	if value == "" {
		return fallback, nil
	}

	parsed, err := time.ParseDuration(value)
	if err != nil {
		return 0, fmt.Errorf("valor invalido para %s: %w", key, err)
	}

	return parsed, nil
}
