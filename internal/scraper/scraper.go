package scraper

import (
	"context"
	"fmt"
	"io"
	"mime"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/gmorais-dev/cadProd/internal/model"
)

const (
	defaultRequestTimeout = 5 * time.Second
	maxHTMLBodySize       = 2 << 20
)

var pricePattern = regexp.MustCompile(`(?i)(R\$\s?\d{1,3}(?:\.\d{3})*,\d{2}|\$\s?\d{1,3}(?:,\d{3})*(?:\.\d{2})?)`)

type Scraper struct {
	client *http.Client
}

func NewScraper(requestTimeout ...time.Duration) *Scraper {
	timeout := defaultRequestTimeout

	if len(requestTimeout) > 0 && requestTimeout[0] > 0 {
		timeout = requestTimeout[0]
	}

	return NewScraperWithClient(&http.Client{Timeout: timeout})
}

func NewScraperWithClient(client *http.Client) *Scraper {
	if client == nil {
		client = &http.Client{Timeout: defaultRequestTimeout}
	}

	if client.Timeout <= 0 {
		client.Timeout = defaultRequestTimeout
	}

	return &Scraper{client: client}
}

func (s *Scraper) Scrape(ctx context.Context, job model.ScrapeJob) model.ScrapeResult {
	start := time.Now()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, job.URL, nil)
	if err != nil {
		return failedResult(job, 0, err.Error(), start)
	}
	req.Header.Set("User-Agent", "cadProd-scraper/1.0")
	req.Header.Set("Accept", "text/html,application/xhtml+xml")

	resp, err := s.client.Do(req)
	if err != nil {
		return failedResult(job, 0, err.Error(), start)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return failedResult(job, resp.StatusCode, fmt.Sprintf("status code invalido: %d", resp.StatusCode), start)
	}

	if !isSupportedContentType(resp.Header.Get("Content-Type")) {
		return failedResult(job, resp.StatusCode, fmt.Sprintf("content type nao suportado: %s", resp.Header.Get("Content-Type")), start)
	}

	doc, err := goquery.NewDocumentFromReader(io.LimitReader(resp.Body, maxHTMLBodySize))
	if err != nil {
		return failedResult(job, resp.StatusCode, err.Error(), start)
	}

	price := extractPrice(doc)

	if price == "" {
		return failedResult(job, resp.StatusCode, "preco nao encontrado", start)
	}

	return model.ScrapeResult{
		JobID:      job.ID,
		URL:        job.URL,
		HTTPStatus: resp.StatusCode,
		Price:      price,
		Success:    true,
		Duration:   time.Since(start),
	}
}

func extractPrice(doc *goquery.Document) string {
	attributeSelectors := []struct {
		selector string
		attr     string
	}{
		{selector: "meta[itemprop='price']", attr: "content"},
		{selector: "meta[property='product:price:amount']", attr: "content"},
		{selector: "[data-price]", attr: "data-price"},
	}

	for _, selector := range attributeSelectors {
		value, exists := doc.Find(selector.selector).First().Attr(selector.attr)
		if !exists {
			continue
		}

		if price := normalizePrice(value); price != "" {
			return price
		}
	}

	textSelectors := []string{
		".price",
		".product-price",
		".sale-price",
		".price-current",
		"[data-testid='price']",
		"[class*='price']",
		".a-price .a-offscreen",
	}

	for _, selector := range textSelectors {
		text := normalizeText(doc.Find(selector).First().Text())
		if price := normalizePrice(text); price != "" {
			return price
		}
	}

	return normalizePrice(doc.Text())
}

func normalizePrice(text string) string {
	normalized := normalizeText(text)
	if normalized == "" {
		return ""
	}

	return pricePattern.FindString(normalized)
}

func normalizeText(text string) string {
	return strings.Join(strings.Fields(text), " ")
}

func isSupportedContentType(contentType string) bool {
	if contentType == "" {
		return true
	}

	mediaType, _, err := mime.ParseMediaType(contentType)
	if err != nil {
		return false
	}

	if mediaType == "text/html" || mediaType == "application/xhtml+xml" {
		return true
	}

	return strings.HasPrefix(mediaType, "text/")
}

func failedResult(job model.ScrapeJob, httpStatus int, message string, start time.Time) model.ScrapeResult {
	return model.ScrapeResult{
		JobID:        job.ID,
		URL:          job.URL,
		HTTPStatus:   httpStatus,
		Success:      false,
		ErrorMessage: message,
		Duration:     time.Since(start),
	}
}
