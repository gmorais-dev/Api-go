package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/gmorais-dev/cadProd/internal/config"
	"github.com/gmorais-dev/cadProd/internal/model"
	"github.com/gmorais-dev/cadProd/internal/scraper"
	"github.com/gmorais-dev/cadProd/internal/storage"
	"github.com/gmorais-dev/cadProd/internal/worker"
)

func main() {
	signalCtx, stop := signal.NotifyContext(
		context.Background(),
		os.Interrupt,
		syscall.SIGTERM,
	)
	defer stop()

	if err := run(signalCtx); err != nil {
		log.Fatal(err)
	}
}

func run(signalCtx context.Context) error {
	appConfig, err := config.NewAppConfig()
	if err != nil {
		return err
	}

	urls, err := config.LoadURLs(appConfig.InputFilePath)
	if err != nil {
		return err
	}

	jobs := buildJobs(urls)

	s := scraper.NewScraper(appConfig.JobTimeout)

	globalCtx, cancel := context.WithTimeout(signalCtx, appConfig.GlobalTimeout)
	defer cancel()

	results := worker.RunWorkerPool(
		globalCtx,
		jobs,
		appConfig.WorkerCount,
		appConfig.JobTimeout,
		s,
	)

	err = storage.SaveResultsToCSV(appConfig.OutputFilePath, results)
	if err != nil {
		return err
	}

	successCount := countSuccessfulResults(results)
	failureCount := len(results) - successCount

	fmt.Printf(
		"Resultados salvos em %s | sucessos: %d | falhas: %d\n",
		appConfig.OutputFilePath,
		successCount,
		failureCount,
	)

	for _, result := range results {
		fmt.Printf(
			"Job %d | URL: %s | HTTP: %d | Preco: %s | Sucesso: %v | Erro: %s | Duracao: %s\n",
			result.JobID,
			result.URL,
			result.HTTPStatus,
			result.Price,
			result.Success,
			result.ErrorMessage,
			result.Duration,
		)
	}

	if signalCtx.Err() != nil {
		fmt.Println("processamento interrompido por sinal do sistema")
	}

	return nil
}

func buildJobs(urls []string) []model.ScrapeJob {
	jobs := make([]model.ScrapeJob, 0, len(urls))

	for i, url := range urls {
		jobs = append(jobs, model.ScrapeJob{
			ID:  i + 1,
			URL: url,
		})
	}

	return jobs
}

func countSuccessfulResults(results []model.ScrapeResult) int {
	successCount := 0

	for _, result := range results {
		if result.Success {
			successCount++
		}
	}

	return successCount
}
