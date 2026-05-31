# --- STAGE 1: Builder ---
FROM golang:1.26.3-alpine3.22 AS builder
WORKDIR /app

# Install ca-certificates and git (diperlukan untuk download module)
RUN apk add --no-cache git ca-certificates

# Set GOPROXY agar lebih stabil saat download di cloud server
ENV GOPROXY=https://proxy.golang.org,direct

COPY go.mod go.sum ./
RUN go mod download

COPY . .
# Batasi penggunaan core CPU untuk mencegah Out-Of-Memory (OOM) di container gratis
RUN CGO_ENABLED=0 GOOS=linux GOMAXPROCS=1 go build -trimpath -ldflags="-w -s" -o app ./cmd/api

# --- STAGE 2: Production ---
FROM alpine:3.22 AS production

RUN apk add --no-cache ca-certificates tzdata

ENV TZ=Asia/Jakarta

RUN adduser -D -g '' appuser
WORKDIR /app

COPY --from=builder --chown=appuser:appuser /app/app .
COPY --from=builder --chown=appuser:appuser /app/migrations ./migrations

USER appuser

EXPOSE 8080
CMD ["./app"]
