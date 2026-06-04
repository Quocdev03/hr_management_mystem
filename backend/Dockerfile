# Stage 1: Build the Go binary
FROM golang:1.26-alpine AS builder

# Install build dependencies
RUN apk add --no-cache git gcc musl-dev

WORKDIR /app

# Copy dependency files and download
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the source code
COPY . .

# Build the server and setup binaries
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o server ./cmd/server/main.go
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o setup ./cmd/setup/main.go

# Stage 2: Final lightweight image
FROM alpine:3.19

RUN apk add --no-cache ca-certificates tzdata netcat-openbsd

WORKDIR /app

# Copy binaries and entrypoint script from builder stage
COPY --from=builder /app/server .
COPY --from=builder /app/setup .
COPY --from=builder /app/docker-entrypoint.sh .

# Grant execute permissions
RUN chmod +x docker-entrypoint.sh server setup

# Expose port
EXPOSE 8080

ENTRYPOINT ["./docker-entrypoint.sh"]
