package config

import (
	"bufio"
	"fmt"
	"net/url"
	"os"
	"strings"
)

func LoadURLs(filePath string) ([]string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var urls []string

	scanner := bufio.NewScanner(file)
	lineNumber := 0

	for scanner.Scan() {
		lineNumber++
		line := strings.TrimSpace(scanner.Text())

		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		parsedURL, err := url.ParseRequestURI(line)
		if err != nil || parsedURL.Scheme == "" || parsedURL.Host == "" {
			return nil, fmt.Errorf("URL invalida na linha %d: %q", lineNumber, line)
		}

		urls = append(urls, line)
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	if len(urls) == 0 {
		return nil, fmt.Errorf("nenhuma URL valida encontrada em %s", filePath)
	}

	return urls, nil
}
