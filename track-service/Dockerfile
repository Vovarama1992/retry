# syntax=docker/dockerfile:1

# Стадия сборки
FROM golang:1.23 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . ./
RUN go build -o track-service ./track-service/cmd/main.go

# Стадия рантайма
FROM debian:bookworm-slim

WORKDIR /app

COPY --from=builder /app/track-service /app/track-service

EXPOSE 8080

CMD ["./track-service"]