# SignalArc Backend - Local Development Dockerfile
# Multi-stage build for a small, maintainable Go image.

FROM golang:1.26-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /signalarc-api ./cmd/api

FROM alpine:3.20

RUN apk add --no-cache ca-certificates

COPY --from=builder /signalarc-api /usr/local/bin/signalarc-api

EXPOSE 4000

ENTRYPOINT ["signalarc-api"]
