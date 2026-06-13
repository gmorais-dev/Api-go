package worker

import (
	"context"
	"testing"
	"time"

	"github.com/gmorais-dev/cadProd/internal/model"
)

type fakeScraper struct {
	delayByJobID map[int]time.Duration
	priceByJobID map[int]string
}

func (f fakeScraper) Scrape(ctx context.Context, job model.ScrapeJob) model.ScrapeResult {
	delay := f.delayByJobID[job.ID]

	select {
	case <-ctx.Done():
		return model.ScrapeResult{
			JobID:        job.ID,
			URL:          job.URL,
			Success:      false,
			ErrorMessage: ctx.Err().Error(),
		}
	case <-time.After(delay):
	}

	price := f.priceByJobID[job.ID]
	if price == "" {
		price = "R$ 199,90"
	}

	return model.ScrapeResult{
		JobID:      job.ID,
		URL:        job.URL,
		HTTPStatus: 200,
		Price:      price,
		Success:    true,
	}
}

func TestRunWorkerPoolPreservesInputOrder(t *testing.T) {
	jobs := []model.ScrapeJob{
		{ID: 1, URL: "http://site-a"},
		{ID: 2, URL: "http://site-b"},
		{ID: 3, URL: "http://site-c"},
	}

	results := RunWorkerPool(
		context.Background(),
		jobs,
		2,
		500*time.Millisecond,
		fakeScraper{
			delayByJobID: map[int]time.Duration{
				1: 80 * time.Millisecond,
				2: 10 * time.Millisecond,
				3: 20 * time.Millisecond,
			},
			priceByJobID: map[int]string{
				1: "R$ 100,00",
				2: "R$ 200,00",
				3: "R$ 300,00",
			},
		},
	)

	if len(results) != len(jobs) {
		t.Fatalf("esperado %d resultados, recebido %d", len(jobs), len(results))
	}

	for index, result := range results {
		if result.JobID != jobs[index].ID {
			t.Fatalf("ordem invalida na posicao %d: esperado job %d recebido %d", index, jobs[index].ID, result.JobID)
		}
	}
}

func TestRunWorkerPoolMarksCanceledJobs(t *testing.T) {
	jobs := []model.ScrapeJob{
		{ID: 1, URL: "http://site-a"},
		{ID: 2, URL: "http://site-b"},
		{ID: 3, URL: "http://site-c"},
	}

	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()

	results := RunWorkerPool(
		ctx,
		jobs,
		1,
		500*time.Millisecond,
		fakeScraper{
			delayByJobID: map[int]time.Duration{
				1: 200 * time.Millisecond,
			},
		},
	)

	if len(results) != len(jobs) {
		t.Fatalf("esperado %d resultados, recebido %d", len(jobs), len(results))
	}

	if results[1].ErrorMessage == "" || results[2].ErrorMessage == "" {
		t.Fatal("jobs nao processados deveriam conter mensagem de erro")
	}
}
