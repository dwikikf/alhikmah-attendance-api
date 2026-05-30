# --- STAGE 1: Base ---
FROM golang:1.26.3-alpine3.22 AS base
WORKDIR /app

RUN apk add --no-cache git ca-certificates

COPY go.mod go.sum ./
RUN go mod download

# --- STAGE 2: Development ---
FROM base AS dev
RUN go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@v4.19.1

COPY . .
CMD ["go", "run", "./cmd/api"]

# --- STAGE 3: Builder ---
FROM base AS builder
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags="-w -s" -o app ./cmd/api

# --- STAGE 4: Production ---
FROM alpine:3.22 AS production

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
