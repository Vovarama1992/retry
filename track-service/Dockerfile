# syntax=docker/dockerfile:1

FROM golang:1.23 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Сборка бинарника в отдельный файл
RUN go build -o track-service-binary ./track-service/cmd/main.go

FROM debian:bookworm-slim

# В финальном образе мы создаём структуру /app/track-service
WORKDIR /app/track-service

# Кладём бинарник
COPY --from=builder /app/track-service-binary ./track-service

# Кладём docs (swagger.json, .go, .yaml)
COPY --from=builder /app/track-service/docs ./docs

EXPOSE 8080

CMD ["./track-service"]