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
Dockerfile
docker-compose.yml
internal/
  config/
  model/
  scraper/
  storage/
  worker/
mock-site/
input/
```

## Como executar localmente

### 1. Subir o mock local

```bash
python3 -m http.server 8080 --directory mock-site
```

### 2. Rodar a aplicacao

```bash
go run .
```

O resultado sera salvo em `output/results.csv`.

## Como executar com Docker

### 1. Build da imagem

```bash
docker build -t api-go-scraper .
```

### 2. Rodar somente a aplicacao em container

Quando usar a imagem diretamente, informe um arquivo de entrada acessivel pelo container e garanta que as URLs apontem para hosts resolviveis dentro do container.

```bash
docker run --rm \
  -v "$(pwd)/output:/app/output" \
  api-go-scraper
```

## Como executar com Docker Compose

O Compose sobe dois servicos:

- `mock-site`: servidor HTTP estatico para os arquivos em `mock-site/`;
- `scraper`: aplicacao Go, configurada para ler `input/urls.docker.txt` e acessar `http://mock-site:8080/produto.html` pela rede interna do Compose.

```bash
docker compose up --build scraper
```

O resultado sera salvo em `output/results.csv` no host por meio do volume configurado no Compose.

Para remover os containers criados:

```bash
docker compose down
```

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
