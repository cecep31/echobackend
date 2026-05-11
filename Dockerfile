# --- Builder Stage ---
FROM golang:1.26 AS builder

WORKDIR /app

# Download dependencies first (layer cache)
COPY go.mod go.sum ./
RUN go mod download

# Copy source and build
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o bin/main cmd/main.go

# --- Final Stage ---
FROM debian:trixie-slim

# Install CA certificates for TLS connections (S3, Valkey, external APIs)
RUN apt-get update \
    && apt-get install -y --no-install-recommends ca-certificates \
    && rm -rf /var/lib/apt/lists/*

# Create a non-root user
RUN groupadd -r appgroup && useradd -r -g appgroup appuser

WORKDIR /app

# Copy binary with correct ownership
COPY --from=builder --chown=appuser:appgroup /app/bin/main .

EXPOSE 8080

USER appuser

# Healthcheck — Fly.io and Docker will use this for liveness detection.
# /health is registered in main.go.
HEALTHCHECK --interval=30s --timeout=5s --start-period=10s --retries=3 \
    CMD ["/bin/sh", "-c", "wget -qO- http://localhost:8080/health || exit 1"]

CMD ["./main"]
