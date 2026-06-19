# ═══════════════════════════════════════════════════════════════════════════
# Backend Dockerfile — Production Multi-Stage Build
# ═══════════════════════════════════════════════════════════════════════════

# Stage 1: Build Go binaries
FROM golang:1.25-alpine AS builder

# Chỉ cần git; không cần gcc/musl-dev vì CGO_ENABLED=0
RUN apk add --no-cache git

WORKDIR /app

# Cache dependencies layer (rebuild chỉ khi go.mod/go.sum thay đổi)
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build 3 binaries: server, migrate, seed
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o server  ./cmd/server/main.go
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o migrate ./cmd/migrate/main.go
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o seed    ./cmd/seed/main.go

# Stage 2: Minimal runtime image
FROM alpine:3.19

# ca-certificates: HTTPS calls
# tzdata: timezone support
# netcat-openbsd: wait-for-service trong entrypoint
# wget: healthcheck (thay curl)
RUN apk add --no-cache ca-certificates tzdata netcat-openbsd wget

WORKDIR /app

# Copy binaries từ builder
COPY --from=builder /app/server  .
COPY --from=builder /app/migrate .
COPY --from=builder /app/seed    .

# Copy entrypoints
COPY docker-entrypoint.sh     .
COPY docker-entrypoint.dev.sh .

# Grant execute permissions
RUN chmod +x docker-entrypoint.sh docker-entrypoint.dev.sh server migrate seed

# Expose port
EXPOSE 8080

# Health check dùng wget (có sẵn trong alpine, không cần curl)
HEALTHCHECK --interval=30s --timeout=10s --retries=3 \
    CMD wget --quiet --tries=1 --spider http://localhost:${APP_PORT:-8080}/api/v1/health || exit 1

ENTRYPOINT ["./docker-entrypoint.sh"]
