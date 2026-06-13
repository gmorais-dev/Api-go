package scraper

import (
	"context"
	"io"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/gmorais-dev/cadProd/internal/model"
)

type roundTripFunc func(*http.Request) (*http.Response, error)

func (fn roundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return fn(req)
}

func TestScrapePrice(t *testing.T) {
	html := `
		<html>
			<body>
				<span class="price">R$ 199,90</span>
			</body>
		</html>
		`

	client := &http.Client{
		Timeout: 2 * time.Second,
		Transport: roundTripFunc(func(req *http.Request) (*http.Response, error) {
			return &http.Response{
				StatusCode: http.StatusOK,
				Header: http.Header{
					"Content-Type": []string{"text/html; charset=utf-8"},
				},
				Body: io.NopCloser(strings.NewReader(html)),
			}, nil
		}),
	}

	s := NewScraperWithClient(client)

	job := model.ScrapeJob{
		ID:  1,
		URL: "https://example.com/produto",
	}

	result := s.Scrape(context.Background(), job)

	if !result.Success {
		t.Fatalf("esperava sucesso, recebeu erro: %s", result.ErrorMessage)
	}

	if result.Price != "R$ 199,90" {
		t.Fatalf(
			"preco incorreto. esperado=%s recebido=%s",
			"R$ 199,90",
			result.Price,
		)
	}
}

func TestScrapePriceFromMetaTag(t *testing.T) {
	html := `
	<html>
		<head>
			<meta itemprop="price" content="R$ 299,90">
		</head>
	</html>
	`

	client := &http.Client{
		Timeout: 2 * time.Second,
		Transport: roundTripFunc(func(req *http.Request) (*http.Response, error) {
			return &http.Response{
				StatusCode: http.StatusOK,
				Header: http.Header{
					"Content-Type": []string{"text/html; charset=utf-8"},
				},
				Body: io.NopCloser(strings.NewReader(html)),
			}, nil
		}),
	}

	s := NewScraperWithClient(client)
	result := s.Scrape(context.Background(), model.ScrapeJob{ID: 1, URL: "https://example.com/meta"})

	if !result.Success {
		t.Fatalf("esperava sucesso, recebeu erro: %s", result.ErrorMessage)
	}

	if result.Price != "R$ 299,90" {
		t.Fatalf("preco incorreto. esperado=R$ 299,90 recebido=%s", result.Price)
	}
}
