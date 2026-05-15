# 🏢 HRM API - Human Resource Management

RESTful API bằng Golang cho hệ thống quản lý nhân sự. **Fresher Backend Developer**.

---

## 📁 Cấu trúc Project

```
backend/
├── cmd/
│   └── server/
│       └── main.go                    # Entry point, khởi tạo và chạy server
├── internal/
│   ├── config/
│   │   ├── config.go                  # Đọc biến môi trường
│   │   ├── database.go                # Kết nối MySQL + Auto Migration
│   │   └── seed.go                    # Dữ liệu mẫu (seed data)
│   ├── common/
│   │   └── common.go                  # Hàm dùng chung: parse ID, phân trang...
│   ├── model/
│   │   ├── base.go                    # TimestampModel, PaginationQuery, PaginatedResult
│   │   ├── user.go                    # User, Role, Permission
│   │   ├── employee.go                # Employee + Request DTOs
│   │   ├── department.go              # Department + Request DTOs
│   │   └── auth.go                    # Login/Register Request, LoginResponse, JWT Claims
│   ├── repository/
│   │   ├── user_repository.go         # Interface + Implementation cho User
│   │   ├── employee_repository.go     # Interface + Implementation cho Employee
│   │   └── department_repository.go   # Interface + Implementation cho Department
│   ├── service/
│   │   ├── auth_service.go            # Interface + Implementation xác thực
│   │   ├── employee_service.go        # Interface + Implementation nhân viên
│   │   └── department_service.go      # Interface + Implementation phòng ban
│   ├── handler/
│   │   ├── auth_handler.go            # HTTP handler: Login, GetProfile
│   │   ├── employee_handler.go        # HTTP handler: CRUD nhân viên
│   │   └── department_handler.go      # HTTP handler: CRUD phòng ban
│   ├── middleware/
│   │   └── middleware.go              # Logger, Recovery, JWT Auth, RBAC
│   └── router/
│       └── router.go                  # Định nghĩa tất cả routes
├── package/
│   ├── response/
│   │   └── response.go                # Unified JSON response format
│   ├── utils/
│   │   └── jwt.go                     # Generate + Validate JWT token
│   ├── validation/
│   │   └── validation.go              # Input validation: email, phone, date...
│   └── loggers/
│       └── logger.go                  # Logger wrapper
├── .env                               # Biến môi trường
├── go.mod
└── go.sum
```

### Kiến trúc phân lớp

```
Request → Router → Middleware → Handler → Service → Repository → Database
                                  ↓            ↓
                              Validation    Validation
                             (handler)    (defense-in-depth)
```

---

## 🚀 Hướng dẫn chạy project

### Cách 1: Docker Compose (Khuyến nghị)

```bash
# 1. Clone project
git clone <repo-url>
cd hrm-api

# 2. Copy env file
cp .env.example .env

# 3. Chạy tất cả services (MySQL + App)
make docker-up
# hoặc: docker-compose up -d --build

# 4. Kiểm tra server
curl http://localhost:8080/health

# 5. Xem logs
make docker-logs
```

### Cách 2: Chạy Local

```bash
# Yêu cầu: Go 1.26+, MySQL 8.0+

# 1. Cài dependencies
go mod tidy

# 2. Tạo database MySQL
mysql -u root -p
CREATE DATABASE hrm_db;

# 3. Copy và điền thông tin .env
cp .env.example .env
# Sửa DB_HOST=localhost, DB_USER=root, DB_PASSWORD=...

# 4. Chạy server
make run
# hoặc: go run ./cmd/server/main.go
```

---

## 🔑 Tài khoản mẫu (Seed Data)

| Username     | Password      | Role                     |
| ------------ | ------------- | ------------------------ |
| `admin`      | `password123` | Admin - toàn quyền       |
| `hr_manager` | `password123` | HR - quản lý nhân viên   |
| `john_doe`   | `password123` | Employee - xem thông tin |

---

## 📖 API Documentation

### Base URL

```
http://localhost:8080/api/v1
```

### Response Format (thống nhất)

**Thành công:**

```json
{
  "success": true,
  "message": "Thao tác thành công",
  "data": {}
}
```

**Lỗi validation (422):**

```json
{
  "success": false,
  "message": "Dữ liệu không hợp lệ",
  "error": [
    { "field": "email", "message": "Email không đúng định dạng" },
    { "field": "phone", "message": "Số điện thoại phải bắt đầu bằng 0 và có đúng 10 số" }
  ]
}
```

**Lỗi thông thường (400/401/404/500):**

```json
{
  "success": false,
  "message": "Mô tả lỗi"
}
```

---

### 🔐 Auth APIs

#### Đăng nhập

```http
POST /api/v1/auth/login
Content-Type: application/json

{
  "email": "admin@hrm.com",
  "password": "password123"
}
```

Response:

```json
{
  "success": true,
  "message": "Đăng nhập thành công",
  "data": {
    "access_token": "eyJhbGciOiJIUzI1NiIs...",
    "user": { ... }
  }
}
```

#### Xem profile (cần token)

```http
GET /api/v1/auth/profile
Authorization: Bearer <token>
```

---

### 👥 Employee APIs (cần token)

#### Lấy danh sách nhân viên

```http
GET /api/v1/employees?page=1&limit=10&search=nguyen
Authorization: Bearer <token>
```

> Phân trang tự động chuẩn hoá: page < 1 → 1, limit < 1 → 10, limit > 100 → 100

#### Lấy thông tin 1 nhân viên

```http
GET /api/v1/employees/1
Authorization: Bearer <token>
```

#### Tạo nhân viên mới (role: admin, hr)

```http
POST /api/v1/employees
Authorization: Bearer <token>
Content-Type: application/json

{
  "department_id": 1,
  "first_name": "Nguyễn",
  "last_name": "Văn D",
  "email": "nguyenvand@hrm.com",
  "phone": "0901234570",
  "position": "Junior Developer",
  "salary": 12000000,
  "join_date": "2024-01-15"
}
```

**Validation rules:**

| Field           | Rule                                                   |
| --------------- | ------------------------------------------------------ |
| `department_id` | Bắt buộc, phải tồn tại trong DB                       |
| `first_name`    | Bắt buộc, 2-100 ký tự, chỉ chữ cái + dấu cách        |
| `last_name`     | Bắt buộc, 2-100 ký tự, chỉ chữ cái + dấu cách        |
| `email`         | Bắt buộc, đúng format email, không trùng               |
| `phone`         | Bắt buộc, bắt đầu bằng 0, đúng 10 số                  |
| `position`      | Tùy chọn, nếu có: 2-100 ký tự                         |
| `salary`        | Không được âm                                          |
| `join_date`     | Tùy chọn, format YYYY-MM-DD, không được ở tương lai   |

#### Cập nhật nhân viên (role: admin, hr)

```http
PUT /api/v1/employees/1
Authorization: Bearer <token>
Content-Type: application/json

{
  "position": "Senior Developer",
  "salary": 20000000,
  "status": "active"
}
```

> `status` chỉ chấp nhận: `active` hoặc `inactive`

#### Xóa nhân viên (role: admin only)

```http
DELETE /api/v1/employees/1
Authorization: Bearer <token>
```

---

### 🏢 Department APIs (cần token)

#### Lấy danh sách phòng ban

```http
GET /api/v1/departments?page=1&limit=10&search=IT
Authorization: Bearer <token>
```

#### Lấy thông tin 1 phòng ban

```http
GET /api/v1/departments/1
Authorization: Bearer <token>
```

#### Tạo phòng ban mới (role: admin)

```http
POST /api/v1/departments
Authorization: Bearer <token>
Content-Type: application/json

{
  "name": "Marketing",
  "code": "MKT",
  "description": "Phòng marketing và truyền thông"
}
```

**Validation rules:**

| Field         | Rule                                                              |
| ------------- | ----------------------------------------------------------------- |
| `name`        | Bắt buộc, 1-100 ký tự, không trùng                               |
| `code`        | Bắt buộc, 1-20 ký tự, chỉ chữ/số/gạch, tự động uppercase        |
| `description` | Tùy chọn, tối đa 500 ký tự                                       |

#### Cập nhật phòng ban (role: admin)

```http
PUT /api/v1/departments/1
Authorization: Bearer <token>
Content-Type: application/json

{
  "name": "IT Department",
  "description": "Updated description",
  "manager_id": 5
}
```

> `manager_id` phải là nhân viên thuộc phòng ban đó

#### Xóa phòng ban (role: admin)

```http
DELETE /api/v1/departments/1
Authorization: Bearer <token>
```

> Không thể xóa phòng ban còn nhân viên

---

## 🛡️ Input Validation

Validation được thực hiện **2 lớp** (defense-in-depth):

1. **Handler Layer** — Trim whitespace, normalize (email lowercase, code uppercase), validate format bằng `package/validation`, trả lỗi 422 chi tiết theo field
2. **Service Layer** — Re-validate, kiểm tra business logic (email trùng, phòng ban tồn tại, ngày không ở tương lai...)

Tất cả input string được tự động **trim whitespace** trước khi xử lý.

---

## 🧠 Các kỹ thuật Golang được minh họa

### 1. Context (Graceful Shutdown)

```go
// cmd/server/main.go
ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
defer cancel()
srv.Shutdown(ctx)
```

### 2. Interface + Dependency Injection

```go
// Định nghĩa interface (repository/user_repository.go)
type UserRepository interface {
    Create(user *model.User) error
    FindByID(id uint) (*model.User, error)
}
// Implement interface
type userRepository struct { db *gorm.DB }
// Inject vào service
authSvc := service.NewAuthService(userRepo, &cfg.JWT)
```

### 3. Tách file theo Domain

```go
// Tất cả file trong cùng folder = cùng package
// → Không cần thay đổi import khi tách file
handler/
├── auth_handler.go          // package handler
├── employee_handler.go      // package handler
└── department_handler.go    // package handler
```

---

## 🔒 Phân quyền RBAC

| Action         | Employee | HR  | Admin |
| -------------- | -------- | --- | ----- |
| Xem nhân viên  | ✅       | ✅  | ✅    |
| Tạo nhân viên  | ❌       | ✅  | ✅    |
| Sửa nhân viên  | ❌       | ✅  | ✅    |
| Xóa nhân viên  | ❌       | ❌  | ✅    |
| Xem phòng ban  | ✅       | ✅  | ✅    |
| CRUD phòng ban | ❌       | ❌  | ✅    |

---

## 🛠 Tech Stack

| Thành phần | Công nghệ               |
| ---------- | ----------------------- |
| Language   | Go 1.26                 |
| Framework  | Gin                     |
| ORM        | GORM                    |
| Database   | MySQL 8.0               |
| Auth       | JWT (golang-jwt/jwt v5) |
| Container  | Docker + Docker Compose |

---

## 📝 Interview Notes

**Câu hỏi thường gặp:**

1. **Goroutine vs Thread?** - Goroutine nhẹ hơn, do Go runtime quản lý, có thể chạy hàng triệu goroutine
2. **Channel để làm gì?** - Giao tiếp an toàn giữa các goroutines (don't communicate by sharing memory)
3. **Interface trong Go?** - Implicit implementation, giúp decouple + testable
4. **Soft Delete là gì?** - Set deleted_at thay vì xóa thật, GORM hỗ trợ qua `gorm.DeletedAt`
5. **JWT hoạt động thế nào?** - Header.Payload.Signature, server verify bằng secret key
6. **Defense-in-depth validation?** - Validate ở cả handler (format) lẫn service (business logic) để bảo vệ 2 lớp
