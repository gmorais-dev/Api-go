package model

import "time"

type ScrapeJob struct {
	ID  int
	URL string
}

type ScrapeResult struct {
	JobID        int
	URL          string
	HTTPStatus   int
	Price        string
	Success      bool
	ErrorMessage string
	Duration     time.Duration
}
