package storage

import (
	"os"
	"strings"
	"testing"
	"time"

	"github.com/gmorais-dev/cadProd/internal/model"
)

func TestSaveResultsToCSV(t *testing.T) {
	dir, err := os.MkdirTemp("", "results-dir-*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(dir)

	filePath := dir + "/nested/results.csv"

	results := []model.ScrapeResult{
		{
			JobID:        1,
			URL:          "http://localhost:8080/produto.html",
			HTTPStatus:   200,
			Price:        "R$ 199,90",
			Success:      true,
			ErrorMessage: "",
			Duration:     10 * time.Millisecond,
		},
	}

	err = SaveResultsToCSV(filePath, results)
	if err != nil {
		t.Fatal(err)
	}

	content, err := os.ReadFile(filePath)
	if err != nil {
		t.Fatal(err)
	}

	text := string(content)

	if !strings.Contains(text, "job_id,url,http_status,price,success,error_message,duration_ms") {
		t.Fatal("cabecalho CSV nao encontrado")
	}

	if !strings.Contains(text, "http://localhost:8080/produto.html") {
		t.Fatal("URL nao encontrada no CSV")
	}

	if !strings.Contains(text, "R$ 199,90") {
		t.Fatal("preco nao encontrado no CSV")
	}

	if !strings.Contains(text, "200") {
		t.Fatal("http_status nao encontrado no CSV")
	}

	if !strings.Contains(text, "true") {
		t.Fatal("status success nao encontrado no CSV")
	}
}
