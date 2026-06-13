# Api-go

Orquestrador de web scraping concorrente em Go usando `worker pool`.

## O que o projeto demonstra

- leitura de uma lista de URLs;
- processamento concorrente com goroutines e channels;
- limite controlado de workers;
- timeout por job e timeout global com `context.Context`;
- graceful shutdown com `Ctrl+C`;
- parser HTML com `goquery` para tentar localizar preco;
- persistencia simples em CSV;
- testes para config, scraper, storage e worker pool.

## Estrutura

```text
main.go
internal/
  config/
  model/
  scraper/
  storage/
  worker/
mock-site/
input/
```

## Como executar

### 1. Subir o mock local

```bash
python3 -m http.server 8080 --directory mock-site
```

### 2. Rodar a aplicacao

```bash
go run .
```

O resultado sera salvo em `output/results.csv`.

## Configuracao por ambiente

- `SCRAPER_INPUT_FILE`
- `SCRAPER_OUTPUT_FILE`
- `SCRAPER_WORKER_COUNT`
- `SCRAPER_JOB_TIMEOUT`
- `SCRAPER_GLOBAL_TIMEOUT`

Exemplo:

```bash
SCRAPER_WORKER_COUNT=5 SCRAPER_JOB_TIMEOUT=3s go run .
```

## Testes

```bash
go test ./...
go vet ./...
```
