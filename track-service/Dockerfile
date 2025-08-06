# syntax=docker/dockerfile:1

FROM golang:1.23 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o track-service-binary ./track-service/cmd/main.go

FROM debian:bookworm-slim

WORKDIR /app

COPY --from=builder /app/track-service-binary /app/track-service

EXPOSE 8080

CMD ["./track-service"]