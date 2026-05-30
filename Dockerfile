# --- STAGE 1: Base ---
FROM golang:1.26-alpine AS base
WORKDIR /app

RUN apk add --no-cache git ca-certificates

COPY go.mod go.sum ./
RUN go mod download

# --- STAGE 2: Development ---
FROM base AS dev
RUN go install github.com/air-verse/air@latest && \
    go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

COPY . .
CMD ["air"]

# --- STAGE 3: Builder ---
FROM base AS builder
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags="-w -s" -o app ./cmd/api

# --- STAGE 4: Production ---
FROM alpine:3.20 AS production

# Tambahkan tzdata agar kontainer bisa membaca zona waktu WIB/WITA/WIT
RUN apk add --no-cache ca-certificates tzdata

# Set zona waktu default ke Waktu Indonesia Barat (WIB)
ENV TZ=Asia/Jakarta

RUN adduser -D -g '' appuser
WORKDIR /app

COPY --from=builder --chown=appuser:appuser /app/app .
COPY --from=builder --chown=appuser:appuser /app/migrations ./migrations

USER appuser

EXPOSE 8080
CMD ["./app"]
