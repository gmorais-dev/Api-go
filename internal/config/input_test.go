package config

import (
	"os"
	"testing"
)

func TestLoadURLs(t *testing.T) {
	content := "# comentario\nhttps://site-a.com\n\n https://site-b.com \n"

	file, err := os.CreateTemp("", "urls-*.txt")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(file.Name())

	_, err = file.WriteString(content)
	if err != nil {
		t.Fatal(err)
	}

	err = file.Close()
	if err != nil {
		t.Fatal(err)
	}

	urls, err := LoadURLs(file.Name())
	if err != nil {
		t.Fatal(err)
	}

	if len(urls) != 2 {
		t.Fatalf("esperado 2 URLs, recebido %d", len(urls))
	}

	if urls[0] != "https://site-a.com" {
		t.Fatalf("URL 1 invalida: %s", urls[0])
	}

	if urls[1] != "https://site-b.com" {
		t.Fatalf("URL 2 invalida: %s", urls[1])
	}
}

func TestLoadURLsInvalidLine(t *testing.T) {
	file, err := os.CreateTemp("", "urls-invalid-*.txt")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(file.Name())

	_, err = file.WriteString("https://site-a.com\ninvalida\n")
	if err != nil {
		t.Fatal(err)
	}

	if err := file.Close(); err != nil {
		t.Fatal(err)
	}

	_, err = LoadURLs(file.Name())
	if err == nil {
		t.Fatal("esperava erro para URL invalida")
	}
}
