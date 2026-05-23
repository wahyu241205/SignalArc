# SignalArc Backend - Local Development / Cloud Run Dockerfile
# Multi-stage build for a small, maintainable Go image.
#
# Runtime stage also installs Node.js, npm, and the global @circle-fin/cli
# package so the backend can invoke the Circle Agent Wallet CLI in the same
# way as the local WSL development environment.
#
# This image must NOT bake any Circle credentials, OTP request IDs, session
# files, private keys, or API keys. Circle CLI session strategy in production
# remains an unsolved runtime concern handled outside the image.

FROM golang:1.26-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /signalarc-api ./cmd/api

FROM alpine:3.20

# Runtime dependencies:
# - ca-certificates: required for outbound HTTPS calls
# - nodejs, npm: required to install and run the global Circle CLI
RUN apk add --no-cache \
        ca-certificates \
        nodejs \
        npm \
    && npm install -g @circle-fin/cli \
    && npm cache clean --force

COPY --from=builder /signalarc-api /usr/local/bin/signalarc-api

EXPOSE 4000

ENTRYPOINT ["signalarc-api"]
