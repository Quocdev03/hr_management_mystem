#!/bin/sh

# ═══════════════════════════════════════════════════════════════════════════
# Backend Development Entrypoint — Hot Reload (air)
# ═══════════════════════════════════════════════════════════════════════════
# Chỉ chờ DB + Redis, sau đó start air (hot-reload).
# Migration và Seed KHÔNG chạy tự động.
#
# Chạy migration: docker compose -f docker-compose.dev.yml run --rm migrate
# Chạy seed:      docker compose -f docker-compose.dev.yml run --rm seed
# ═══════════════════════════════════════════════════════════════════════════

set -e

DB_HOST=${DB_HOST:-mysql}
DB_PORT=${DB_PORT:-3306}
REDIS_HOST=${REDIS_HOST:-redis}
REDIS_PORT=${REDIS_PORT:-6379}
WAIT_TIMEOUT=${WAIT_TIMEOUT:-60}

log_info()    { echo "[INFO]    $*"; }
log_success() { echo "[SUCCESS] $*"; }
log_error()   { echo "[ERROR]   $*" >&2; }

# ─── Wait for dependency ──────────────────────────────────────────────────────
wait_for() {
    local host=$1
    local port=$2
    local name=$3
    local elapsed=0

    log_info "Connecting to $name at $host:$port..."
    while ! nc -z "$host" "$port" 2>/dev/null; do
        if [ "$elapsed" -ge "$WAIT_TIMEOUT" ]; then
            log_error "Timeout: $name ($host:$port) not ready after ${WAIT_TIMEOUT}s"
            exit 1
        fi
        sleep 1
        elapsed=$((elapsed + 1))
    done
    log_success "$name connected (${elapsed}s)"
}

# ─── Main ────────────────────────────────────────────────────────────────────
log_info "Loading environment... (APP_ENV=${APP_ENV:-development})"

wait_for "$DB_HOST"    "$DB_PORT"    "MySQL"
wait_for "$REDIS_HOST" "$REDIS_PORT" "Redis"

log_info "Starting hot-reload with air..."
exec air -c .air.toml 2>/dev/null || exec air
