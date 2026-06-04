# ==============================================================================
# HR Management System - Makefile
# ==============================================================================
# Sử dụng: make <target>
# Ví dụ  : make run | make docker-up | make build

# ─── Biến cấu hình ──────────────────────────────────────────────────────────
BINARY_NAME    := hrm-server
BACKEND_DIR    := backend
FRONTEND_DIR   := frontend
SERVER_ENTRY   := $(BACKEND_DIR)/cmd/server/main.go
SETUP_ENTRY    := $(BACKEND_DIR)/cmd/setup/main.go
BUILD_OUT      := ./bin/$(BINARY_NAME)

# Màu terminal (ANSI)
GREEN  := \033[0;32m
YELLOW := \033[0;33m
CYAN   := \033[0;36m
RESET  := \033[0m

.PHONY: help \
        docker-up docker-down docker-logs docker-ps docker-restart \
        run build clean tidy test \
        setup \
        frontend-install frontend-dev frontend-build \
        lint fmt \
        all

# ─── Default target ──────────────────────────────────────────────────────────
.DEFAULT_GOAL := help

## help: Hiển thị danh sách tất cả các lệnh có sẵn
help:
	@echo ""
	@echo "$(CYAN)╔══════════════════════════════════════════════════╗$(RESET)"
	@echo "$(CYAN)║       HR Management System — Makefile Help       ║$(RESET)"
	@echo "$(CYAN)╚══════════════════════════════════════════════════╝$(RESET)"
	@echo ""
	@echo "$(YELLOW)🐳  Docker Commands:$(RESET)"
	@echo "  make docker-up       Khởi chạy container MySQL & Redis (detached)"
	@echo "  make docker-down     Dừng và xoá tất cả container"
	@echo "  make docker-restart  Restart toàn bộ Docker stack"
	@echo "  make docker-logs     Xem log từ container 'app'"
	@echo "  make docker-ps       Liệt kê trạng thái các service"
	@echo ""
	@echo "$(YELLOW)⚙️   Backend (Go) Commands:$(RESET)"
	@echo "  make setup           Khởi tạo DB, migration & seed dữ liệu ban đầu"
	@echo "  make run             Chạy API server (hot-reload sẵn sàng)"
	@echo "  make build           Biên dịch Go → $(BUILD_OUT)"
	@echo "  make clean           Xoá binary đã build"
	@echo "  make tidy            Dọn dẹp & cập nhật Go modules"
	@echo "  make test            Chạy toàn bộ unit test"
	@echo "  make fmt             Format code Go (gofmt)"
	@echo "  make lint            Kiểm tra linting Go (golangci-lint)"
	@echo ""
	@echo "$(YELLOW)🖥️   Frontend (Vue 3) Commands:$(RESET)"
	@echo "  make frontend-install  Cài đặt npm dependencies"
	@echo "  make frontend-dev      Khởi chạy Vite dev server"
	@echo "  make frontend-build    Build production bundle"
	@echo ""
	@echo "$(YELLOW)🚀  Shortcuts:$(RESET)"
	@echo "  make all             Setup → docker-up → run (full start)"
	@echo ""

# ─── Docker ──────────────────────────────────────────────────────────────────

## docker-up: Khởi chạy MySQL & Redis containers ở chế độ nền
docker-up:
	@echo "$(GREEN)▶ Khởi chạy Docker containers...$(RESET)"
	docker-compose up -d
	@echo "$(GREEN)✔ Containers đang chạy. MySQL:3306 | Redis:6379$(RESET)"

## docker-down: Dừng và xoá tất cả containers
docker-down:
	@echo "$(YELLOW)■ Dừng Docker containers...$(RESET)"
	docker-compose down
	@echo "$(GREEN)✔ Đã dừng tất cả containers.$(RESET)"

## docker-restart: Restart toàn bộ Docker stack
docker-restart: docker-down docker-up

## docker-logs: Theo dõi log container 'app'
docker-logs:
	docker-compose logs -f app

## docker-ps: Xem trạng thái các service
docker-ps:
	docker-compose ps

# ─── Backend (Go) ────────────────────────────────────────────────────────────

## setup: Khởi tạo DB schema, migration & seed dữ liệu mẫu
setup:
	@echo "$(GREEN)▶ Đang thiết lập cơ sở dữ liệu...$(RESET)"
	go run $(SETUP_ENTRY)
	@echo "$(GREEN)✔ Thiết lập DB hoàn tất.$(RESET)"

## run: Chạy API server Go trực tiếp (development mode)
run:
	@echo "$(GREEN)▶ Khởi chạy API server tại http://localhost:8080 ...$(RESET)"
	go run $(SERVER_ENTRY)

## build: Biên dịch file thực thi Go
build:
	@echo "$(GREEN)▶ Đang build binary → $(BUILD_OUT)$(RESET)"
	@mkdir -p bin
	go build -ldflags="-s -w" -o $(BUILD_OUT) $(SERVER_ENTRY)
	@echo "$(GREEN)✔ Build hoàn tất: $(BUILD_OUT)$(RESET)"

## clean: Xoá binary đã build
clean:
	@echo "$(YELLOW)■ Xoá $(BUILD_OUT)...$(RESET)"
	@rm -f $(BUILD_OUT)
	@echo "$(GREEN)✔ Đã dọn dẹp.$(RESET)"

## tidy: Dọn dẹp và đồng bộ Go modules
tidy:
	@echo "$(GREEN)▶ Đang chạy go mod tidy...$(RESET)"
	cd $(BACKEND_DIR) && go mod tidy
	@echo "$(GREEN)✔ Go modules đã được cập nhật.$(RESET)"

## test: Chạy toàn bộ test suite
test:
	@echo "$(GREEN)▶ Chạy tests...$(RESET)"
	cd $(BACKEND_DIR) && go test ./... -v -count=1
	@echo "$(GREEN)✔ Tests hoàn tất.$(RESET)"

## fmt: Format toàn bộ code Go
fmt:
	@echo "$(GREEN)▶ Formatting Go code...$(RESET)"
	cd $(BACKEND_DIR) && gofmt -w .
	@echo "$(GREEN)✔ Format hoàn tất.$(RESET)"

## lint: Chạy golangci-lint (cần cài sẵn)
lint:
	@echo "$(GREEN)▶ Chạy golangci-lint...$(RESET)"
	cd $(BACKEND_DIR) && golangci-lint run ./...

# ─── Frontend (Vue 3 + Vite) ─────────────────────────────────────────────────

## frontend-install: Cài đặt npm packages
frontend-install:
	@echo "$(GREEN)▶ Cài đặt frontend dependencies...$(RESET)"
	cd $(FRONTEND_DIR) && npm install
	@echo "$(GREEN)✔ npm install hoàn tất.$(RESET)"

## frontend-dev: Khởi chạy Vite dev server
frontend-dev:
	@echo "$(GREEN)▶ Khởi chạy Vite dev server tại http://localhost:3000 ...$(RESET)"
	cd $(FRONTEND_DIR) && npm run dev

## frontend-build: Build production bundle
frontend-build:
	@echo "$(GREEN)▶ Build frontend cho production...$(RESET)"
	cd $(FRONTEND_DIR) && npm run build
	@echo "$(GREEN)✔ Frontend build hoàn tất → $(FRONTEND_DIR)/dist$(RESET)"

# ─── Shortcut tổng hợp ───────────────────────────────────────────────────────

## all: Chạy toàn bộ: docker-up → setup → run
all: docker-up setup run
