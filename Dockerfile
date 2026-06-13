# syntax=docker/dockerfile:1

FROM golang:1.25-alpine AS builder

WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -trimpath -ldflags="-s -w" -o /bin/scraper .

FROM alpine:3.22

RUN addgroup -S app && adduser -S app -G app
WORKDIR /app

COPY --from=builder /bin/scraper /app/scraper
COPY input /app/input

RUN mkdir -p /app/output && chown -R app:app /app

USER app

ENV SCRAPER_INPUT_FILE=input/urls.txt \
    SCRAPER_OUTPUT_FILE=output/results.csv

ENTRYPOINT ["/app/scraper"]
