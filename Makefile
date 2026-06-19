# ==============================================================================
# HR Management System — Makefile
# ==============================================================================
# Sử dụng: make <target>
# Ví dụ  : make dev | make dev-migrate | make prod

# ─── Biến cấu hình ──────────────────────────────────────────────────────────
BINARY_NAME    := hrm-server
BACKEND_DIR    := backend
FRONTEND_DIR   := frontend
SERVER_ENTRY   := ./$(BACKEND_DIR)/cmd/server/main.go
MIGRATE_ENTRY  := ./$(BACKEND_DIR)/cmd/migrate/main.go
SEED_ENTRY     := ./$(BACKEND_DIR)/cmd/seed/main.go
BUILD_OUT      := ./bin/$(BINARY_NAME)

COMPOSE_DEV    := docker compose -f docker-compose.dev.yml
COMPOSE_PROD   := docker compose -f docker-compose.prod.yml

# Màu terminal (ANSI)
GREEN  := \033[0;32m
YELLOW := \033[0;33m
CYAN   := \033[0;36m
RED    := \033[0;31m
RESET  := \033[0m

.PHONY: help \
        dev dev-down dev-logs dev-logs-backend dev-logs-frontend \
        dev-ps dev-restart dev-build \
        dev-migrate dev-seed \
        prod prod-down prod-logs prod-ps prod-restart prod-migrate \
        backend-run backend-migrate backend-seed \
        backend-build backend-clean backend-tidy backend-test \
        frontend-install frontend-dev frontend-build \
        lint fmt validate \
        all

# ─── Default target ──────────────────────────────────────────────────────────
.DEFAULT_GOAL := help

## help: Hiển thị danh sách tất cả các lệnh có sẵn
help:
	@echo ""
	@echo "$(CYAN)╔══════════════════════════════════════════════════════╗$(RESET)"
	@echo "$(CYAN)║        HR Management System — Makefile Help          ║$(RESET)"
	@echo "$(CYAN)╚══════════════════════════════════════════════════════╝$(RESET)"
	@echo ""
	@echo "$(YELLOW)🚀  Development — Startup:$(RESET)"
	@echo "  make dev              Khởi chạy stack dev (hot-reload)"
	@echo "  make dev-down         Dừng stack dev"
	@echo "  make dev-restart      Restart stack dev"
	@echo "  make dev-build        Rebuild image dev"
	@echo "  make dev-ps           Trạng thái service dev"
	@echo ""
	@echo "$(YELLOW)🗄️   Development — Database:$(RESET)"
	@echo "  make dev-migrate      Chạy migration (dev)"
	@echo "  make dev-seed         Chạy seed dữ liệu mẫu (dev)"
	@echo ""
	@echo "$(YELLOW)📋  Development — Logs:$(RESET)"
	@echo "  make dev-logs         Log tất cả service"
	@echo "  make dev-logs-backend Log backend"
	@echo "  make dev-logs-frontend Log frontend"
	@echo ""
	@echo "$(YELLOW)🏭  Production — Startup:$(RESET)"
	@echo "  make validate         Validate config trước khi deploy"
	@echo "  make prod             Khởi chạy stack production"
	@echo "  make prod-down        Dừng stack production"
	@echo "  make prod-restart     Restart stack production"
	@echo "  make prod-logs        Log production"
	@echo "  make prod-ps          Trạng thái service production"
	@echo ""
	@echo "$(YELLOW)🗄️   Production — Database:$(RESET)"
	@echo "  make prod-migrate     Chạy migration (production)"
	@echo ""
	@echo "$(YELLOW)⚙️   Backend (Go) — chạy trực tiếp (local):$(RESET)"
	@echo "  make backend-run      Chạy server (go run)"
	@echo "  make backend-migrate  Chạy migration (go run, dùng backend/.env)"
	@echo "  make backend-seed     Chạy seed (go run, dùng backend/.env)"
	@echo "  make backend-setup    Chạy migrate + seed trong một bước"
	@echo "  make backend-build    Build binary → $(BUILD_OUT)"
	@echo "  make backend-clean    Xoá binary"
	@echo "  make backend-tidy     go mod tidy"
	@echo "  make backend-test     Chạy unit tests"
	@echo "  make fmt              Format Go code"
	@echo "  make lint             golangci-lint"
	@echo ""
	@echo "$(YELLOW)🖥️   Frontend (Vue 3) — chạy trực tiếp (local):$(RESET)"
	@echo "  make frontend-install   npm install"
	@echo "  make frontend-dev       Vite dev server (port 5173)"
	@echo "  make frontend-build     Build production bundle"
	@echo ""

# ─── Development — Startup ───────────────────────────────────────────────────

## dev: Khởi chạy stack dev (hot-reload). Migration/Seed KHÔNG tự chạy.
dev:
	@echo "$(GREEN)▶ Khởi chạy Dev stack...$(RESET)"
	@echo "$(CYAN)  Frontend : http://localhost:5173$(RESET)"
	@echo "$(CYAN)  Backend  : http://localhost:8080$(RESET)"
	@echo "$(CYAN)  MySQL    : localhost:3306$(RESET)"
	@echo "$(CYAN)  Redis    : localhost:6379$(RESET)"
	$(COMPOSE_DEV) up -d
	@echo "$(GREEN)✔ Dev stack đang chạy.$(RESET)"
	@echo "$(YELLOW)  → Chạy migration lần đầu: make dev-migrate$(RESET)"
	@echo "$(YELLOW)  → Chạy seed dữ liệu mẫu: make dev-seed$(RESET)"

## dev-down: Dừng stack dev
dev-down:
	@echo "$(YELLOW)■ Dừng Dev stack...$(RESET)"
	$(COMPOSE_DEV) down
	@echo "$(GREEN)✔ Dev stack đã dừng.$(RESET)"

## dev-restart: Restart stack dev
dev-restart: dev-down dev

## dev-build: Rebuild image dev
dev-build:
	@echo "$(GREEN)▶ Rebuild Dev images...$(RESET)"
	$(COMPOSE_DEV) build --no-cache
	@echo "$(GREEN)✔ Rebuild hoàn tất.$(RESET)"

## dev-ps: Trạng thái service dev
dev-ps:
	$(COMPOSE_DEV) ps

# ─── Development — Logs ──────────────────────────────────────────────────────

## dev-logs: Log tất cả service
dev-logs:
	$(COMPOSE_DEV) logs -f

## dev-logs-backend: Log backend
dev-logs-backend:
	$(COMPOSE_DEV) logs -f backend

## dev-logs-frontend: Log frontend
dev-logs-frontend:
	$(COMPOSE_DEV) logs -f frontend

# ─── Development — Database ──────────────────────────────────────────────────

## dev-migrate: Chạy migration (Docker, dev). Idempotent — an toàn chạy nhiều lần.
dev-migrate:
	@echo "$(GREEN)▶ Chạy migration (dev)...$(RESET)"
	$(COMPOSE_DEV) run --rm migrate
	@echo "$(GREEN)✔ Migration hoàn tất.$(RESET)"

## dev-seed: Chạy seed dữ liệu mẫu (Docker, dev). Idempotent — không tạo trùng.
dev-seed:
	@echo "$(GREEN)▶ Chạy seed (dev)...$(RESET)"
	$(COMPOSE_DEV) run --rm seed
	@echo "$(GREEN)✔ Seed hoàn tất.$(RESET)"

# ─── Production — Startup ────────────────────────────────────────────────────

## validate: Validate config trước khi deploy production
validate:
	@echo "$(GREEN)▶ Validating production config...$(RESET)"
	./scripts/validate-production.sh

## prod: Khởi chạy stack production. Migration KHÔNG tự chạy.
prod:
	@echo "$(GREEN)▶ Khởi chạy Production stack...$(RESET)"
	$(COMPOSE_PROD) up -d
	@echo "$(GREEN)✔ Production stack chạy tại http://localhost$(RESET)"
	@echo "$(YELLOW)  → Chạy migration: make prod-migrate$(RESET)"

## prod-down: Dừng stack production
prod-down:
	@echo "$(YELLOW)■ Dừng Production stack...$(RESET)"
	$(COMPOSE_PROD) down
	@echo "$(GREEN)✔ Production stack đã dừng.$(RESET)"

## prod-restart: Restart production stack
prod-restart: prod-down prod

## prod-logs: Log production
prod-logs:
	$(COMPOSE_PROD) logs -f

## prod-ps: Trạng thái service production
prod-ps:
	$(COMPOSE_PROD) ps

# ─── Production — Database ───────────────────────────────────────────────────

## prod-migrate: Chạy migration (Docker, production). Idempotent.
prod-migrate:
	@echo "$(GREEN)▶ Chạy migration (production)...$(RESET)"
	$(COMPOSE_PROD) run --rm migrate
	@echo "$(GREEN)✔ Migration (production) hoàn tất.$(RESET)"

# ─── Backend (Go) — chạy trực tiếp ──────────────────────────────────────────

## backend-run: Chạy API server (go run, dùng backend/.env)
backend-run:
	@echo "$(GREEN)▶ Khởi chạy API server tại http://localhost:8080 ...$(RESET)"
	cd $(BACKEND_DIR) && go run ./cmd/server/main.go

## backend-migrate: Chạy migration local (dùng backend/.env)
backend-migrate:
	@echo "$(GREEN)▶ Chạy migration (local)...$(RESET)"
	cd $(BACKEND_DIR) && go run ./cmd/migrate/main.go
	@echo "$(GREEN)✔ Migration hoàn tất.$(RESET)"

## backend-seed: Chạy seed local (dùng backend/.env)
backend-seed:
	@echo "$(GREEN)▶ Chạy seed (local)...$(RESET)"
	cd $(BACKEND_DIR) && go run ./cmd/seed/main.go
	@echo "$(GREEN)✔ Seed hoàn tất.$(RESET)"
## backend-build: Build binary
backend-build:
	@echo "$(GREEN)▶ Build binary → $(BUILD_OUT)$(RESET)"
	@mkdir -p bin
	cd $(BACKEND_DIR) && go build -ldflags="-s -w" -o ../$(BUILD_OUT) ./cmd/server/main.go
	@echo "$(GREEN)✔ Build hoàn tất: $(BUILD_OUT)$(RESET)"

## backend-clean: Xoá binary
backend-clean:
	@echo "$(YELLOW)■ Xoá $(BUILD_OUT)...$(RESET)"
	@rm -f $(BUILD_OUT)
	@echo "$(GREEN)✔ Đã dọn dẹp.$(RESET)"

## backend-tidy: go mod tidy
backend-tidy:
	@echo "$(GREEN)▶ go mod tidy...$(RESET)"
	cd $(BACKEND_DIR) && go mod tidy
	@echo "$(GREEN)✔ Go modules updated.$(RESET)"

## backend-test: Chạy unit tests
backend-test:
	@echo "$(GREEN)▶ Chạy tests...$(RESET)"
	cd $(BACKEND_DIR) && go test ./... -v -count=1
	@echo "$(GREEN)✔ Tests hoàn tất.$(RESET)"

## fmt: Format Go code
fmt:
	@echo "$(GREEN)▶ Formatting Go code...$(RESET)"
	cd $(BACKEND_DIR) && gofmt -w .
	@echo "$(GREEN)✔ Format hoàn tất.$(RESET)"

## lint: golangci-lint
lint:
	@echo "$(GREEN)▶ Chạy golangci-lint...$(RESET)"
	cd $(BACKEND_DIR) && golangci-lint run ./...

# ─── Frontend (Vue 3) ─────────────────────────────────────────────────────────

## frontend-install: npm install
frontend-install:
	@echo "$(GREEN)▶ npm install...$(RESET)"
	cd $(FRONTEND_DIR) && npm install
	@echo "$(GREEN)✔ npm install hoàn tất.$(RESET)"

## frontend-dev: Vite dev server (port 5173)
frontend-dev:
	@echo "$(GREEN)▶ Vite dev server → http://localhost:5173$(RESET)"
	cd $(FRONTEND_DIR) && npm run dev

## frontend-build: Build production bundle
frontend-build:
	@echo "$(GREEN)▶ Build frontend...$(RESET)"
	cd $(FRONTEND_DIR) && npm run build
	@echo "$(GREEN)✔ Build hoàn tất → $(FRONTEND_DIR)/dist$(RESET)"
