# 🧠 HR MANAGER AI AGENT SYSTEM (CLINE + AGENTS UNIFIED)

---

# 🛑 1. KỶ LUẬT THÉP (GLOBAL GUARDRAILS)

## 1.1 /plan FIRST (MANDATORY GATE)
- KHÔNG được sửa code trước khi:
  1. Trình bày `/plan`
  2. Người dùng APPROVE
- Sau approval mới được phép implement

---

## 1.2 NO DANGEROUS EXECUTION
- Cấm lệnh terminal nguy hiểm, recursive, destructive
- Nếu 2 lần lỗi giống nhau → STOP và hỏi user

---

## 1.3 NO UNAUTHORIZED FILE MODIFICATION
- Cấm sửa nếu không nằm trong scope task
- Cấm chỉnh:
  - backend/.env
  - docker-compose.yml
  - database.sql
  - frontend/node_modules
  - backend/bin

---

## 1.4 SCOPE LOCKING
Chỉ được làm việc trong:
- backend/internal/
- backend/cmd/
- frontend/src/
- frontend/api/
- frontend/router/
- frontend/store/

---

# ⚡ 2. TOKEN OPTIMIZATION RULES

- Chỉ đọc file liên quan trực tiếp task
- Không scan toàn bộ repo
- Không dump full file
- Chỉ output patch / diff
- Ưu tiên logic, không lan man

---

# 🧠 3. PROJECT ARCHITECTURE RULES

## Backend (Go)
- Gin + GORM
- Layered:
  handler → service → repository
- Handler: thin
- Service: business logic
- Repository: DB access

## Frontend (Vue 3 + Vite)
- SPA architecture
- API layer riêng
- Store quản lý state

## Infra
- MySQL + Redis
- JWT auth
- RBAC: admin / hr / employee

---

# 🧩 4. SKILL ROUTING ENGINE

## Backend Go
Use:
- go-reviewer
- golang-patterns
- go-concurrency-patterns
- error-handling

## Database
Use:
- database-reviewer
- mysql-patterns
- database-optimizer
- database-migrations-sql-migrations

## Frontend Vue
Use:
- typescript-reviewer
- vue3-performance
- frontend-security-coder

## Git / Workflow
Use:
- git-workflow

## Debugging / Architecture
Use:
- architect
- code-reviewer

---

# 🧠 5. ROLE SYSTEM

- planner → tạo /plan
- architect → validate structure
- go-reviewer → Go quality
- typescript-reviewer → Vue quality
- database-reviewer → DB + SQL
- frontend-security-coder → XSS, token, DOM
- code-reviewer → tổng review
- git-workflow → commit discipline

---

# 🛠️ 6. SLASH COMMAND SYSTEM

## /plan (MANDATORY)
- phân tích task
- file liên quan
- risk
- steps

## /build-fix
- go test ./...
- go vet ./...
- npm run build
- fix errors

## /code-review
- review diff
- severity report

## /security-scan
- detect secrets
- XSS / SQL injection / token leak

## /perf-audit
- backend + frontend performance review

## /db-migrate
- check schema + migration safety

## /checkpoint
- summary dạng commit message

## /diagnose
- check git status + project health

## /rollback
- revert changes safely

---

# 💾 7. MEMORY SYSTEM

## READ FIRST
- đọc agent_memory.md nếu tồn tại

## WRITE BACK
- sau mỗi stage:
  - task
  - files touched
  - result
  - notes

---

# 🎯 8. EXECUTION FLOW (STRICT)

1. Understand request
2. Check scope
3. Create /plan
4. Wait approval
5. Execute minimal diff
6. Run checks (/build-fix nếu cần)
7. Write checkpoint
8. Update memory

---

# 🔐 9. SECURITY RULES

- Không hardcode secrets
- Không commit API keys
- Không sửa config production
- Không thay đổi migration nếu chưa approve
- Validate all inputs (backend + frontend)

---

# 📌 10. FINAL PRINCIPLE

- Minimal change
- Maximum safety
- Strict architecture
- Token efficient
- Always plan before action