package worker

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/gmorais-dev/cadProd/internal/model"
)

const pendingJobError = "job nao processado"

type JobScraper interface {
	Scrape(ctx context.Context, job model.ScrapeJob) model.ScrapeResult
}

type workItem struct {
	index int
	job   model.ScrapeJob
}

type workResult struct {
	index  int
	result model.ScrapeResult
}

func RunWorkerPool(
	ctx context.Context,
	jobs []model.ScrapeJob,
	workerCount int,
	jobTimeout time.Duration,
	s JobScraper,
) []model.ScrapeResult {
	results := buildPendingResults(jobs)

	if workerCount <= 0 {
		fillPendingResultsWithMessage(results, "worker count invalido")
		return results
	}

	if jobTimeout <= 0 {
		jobTimeout = 2 * time.Second
	}

	jobsChan := make(chan workItem, len(jobs))
	resultsChan := make(chan workResult, len(jobs))

	var wg sync.WaitGroup

	for i := 0; i < workerCount; i++ {
		wg.Add(1)

		go func() {
			defer wg.Done()

			for {
				select {
				case <-ctx.Done():
					return

				case item, ok := <-jobsChan:
					if !ok {
						return
					}

					if ctx.Err() != nil {
						return
					}

					jobCtx, cancel := context.WithTimeout(ctx, jobTimeout)
					result := s.Scrape(jobCtx, item.job)
					cancel()

					resultsChan <- workResult{
						index:  item.index,
						result: result,
					}
				}
			}
		}()
	}

	for index, job := range jobs {
		if ctx.Err() != nil {
			break
		}

		jobsChan <- workItem{
			index: index,
			job:   job,
		}
	}
	close(jobsChan)

	go func() {
		wg.Wait()
		close(resultsChan)
	}()

	for result := range resultsChan {
		results[result.index] = result.result
	}

	if ctx.Err() != nil {
		fillPendingResultsWithMessage(results, fmt.Sprintf("job nao processado: %v", ctx.Err()))
	}

	return results
}

func buildPendingResults(jobs []model.ScrapeJob) []model.ScrapeResult {
	results := make([]model.ScrapeResult, len(jobs))

	for index, job := range jobs {
		results[index] = model.ScrapeResult{
			JobID:        job.ID,
			URL:          job.URL,
			Success:      false,
			ErrorMessage: pendingJobError,
		}
	}

	return results
}

func fillPendingResultsWithMessage(results []model.ScrapeResult, message string) {
	for index, result := range results {
		if result.ErrorMessage != pendingJobError {
			continue
		}

		results[index].ErrorMessage = message
	}
}
