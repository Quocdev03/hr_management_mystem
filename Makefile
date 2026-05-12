# ============================================================
# Makefile - Các lệnh thường dùng trong development
# ============================================================

.PHONY: help run build docker-up docker-down test clean

# Hiển thị help
help:
	@echo "Available commands:"
	@echo "  make run          - Chạy server local"
	@echo "  make build        - Build binary"
	@echo "  make docker-up    - Khởi động với Docker Compose"
	@echo "  make docker-down  - Dừng Docker Compose"
	@echo "  make docker-logs  - Xem logs"
	@echo "  make clean        - Xóa binary"

# Chạy server local (cần có .env và MySQL/Redis đang chạy)
run:
	go run ./cmd/server/main.go

# Build binary
build:
	CGO_ENABLED=0 go build -o hrm-api ./cmd/server/main.go

# Docker Compose commands
docker-up:
	docker-compose up -d --build
	@echo "✅ Services started!"
	@echo "📖 API: http://localhost:8080/api/v1"
	@echo "🔍 Health: http://localhost:8080/health"

docker-down:
	docker-compose down

docker-logs:
	docker-compose logs -f app

docker-restart:
	docker-compose restart app

# Xem logs realtime
logs:
	docker-compose logs -f

# Xóa tất cả (bao gồm volumes)
clean-all:
	docker-compose down -v
	rm -f hrm-api

# Chỉ xóa binary
clean:
	rm -f hrm-api

# Tidy dependencies
tidy:
	go mod tidy

# Format code
fmt:
	go fmt ./...

# Vet code
vet:
	go vet ./...
