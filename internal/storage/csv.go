package storage

import (
	"encoding/csv"
	"os"
	"path/filepath"
	"strconv"

	"github.com/gmorais-dev/cadProd/internal/model"
)

func SaveResultsToCSV(filePath string, results []model.ScrapeResult) error {
	if err := os.MkdirAll(filepath.Dir(filePath), 0o755); err != nil {
		return err
	}

	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)

	header := []string{
		"job_id",
		"url",
		"http_status",
		"price",
		"success",
		"error_message",
		"duration_ms",
	}

	if err := writer.Write(header); err != nil {
		return err
	}

	for _, result := range results {
		row := []string{
			strconv.Itoa(result.JobID),
			result.URL,
			strconv.Itoa(result.HTTPStatus),
			result.Price,
			strconv.FormatBool(result.Success),
			result.ErrorMessage,
			strconv.FormatInt(result.Duration.Milliseconds(), 10),
		}

		if err := writer.Write(row); err != nil {
			return err
		}
	}

	writer.Flush()

	if err := writer.Error(); err != nil {
		return err
	}

	return nil
}
