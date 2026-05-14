# 🏢 HRM API - Human Resource Management

RESTful API bằng Golang cho hệ thống quản lý nhân sự. **Fresher Backend Developer**.

---

## 📁 Cấu trúc Project

```
hrm-api/
├── cmd/
│   └── server/
│       └── main.go              # Entry point, khởi tạo và chạy server
├── internal/
│   ├── config/
│   │   ├── config.go            # Đọc biến môi trường
│   │   ├── database.go          # Kết nối MySQL + Auto Migration
│   │   └── seed.go              # Dữ liệu mẫu (seed data)
│   ├── model/
│   │   └── model.go             # Tất cả structs: Entity + DTO + Request/Response
│   ├── repository/
│   │   └── repository.go        # Interface + Implementation tầng DB
│   ├── service/
│   │   └── service.go           # Business logic + Goroutine + WaitGroup
│   ├── handler/
│   │   └── handler.go           # HTTP handlers (controllers)
│   ├── middleware/
│   │   └── middleware.go        # Logger, Recovery, JWT Auth, RBAC
│   └── router/
│       └── router.go            # Định nghĩa tất cả routes
├── pkg/
│   ├── response/
│   │   └── response.go          # Unified JSON response format
│   ├── utils/
│   │   └── jwt.go               # Generate + Validate JWT token
│   └── logger/
│       └── logger.go            # Logger wrapper
├── configs/                     # File cấu hình bổ sung (nếu cần)
├── .env                         # Biến môi trường
├── .gitignore
├── docker-compose.yml
├── Makefile
└── README.md
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
# Yêu cầu: Go 1.21+, MySQL 8.0+

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

```json
{
	"success": true,
	"message": "Success message",
	"data": {}
}
```

---

### 🔐 Auth APIs

#### Đăng ký tài khoản

```http
POST /api/v1/auth/register
Content-Type: application/json

{
  "username": "newuser",
  "email": "newuser@example.com",
  "password": "password123",
  "role_id": 3
}
```

#### Đăng nhập

```http
POST /api/v1/auth/login
Content-Type: application/json

{
  "username": "admin",
  "password": "password123"
}
```

Response:

```json
{
  "success": true,
  "message": "Login successful",
  "data": {
    "access_token": "eyJhbGciOiJIUzI1NiIs...",
    "refesh_token": "eyJhbGciOiJIUzI1NiIs...",
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

#### Cập nhật nhân viên (role: admin, hr)

```http
PUT /api/v1/employees/1
Authorization: Bearer <token>
Content-Type: application/json

{
  "position": "Senior Developer",
  "salary": 20000000
}
```

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

#### Cập nhật phòng ban (role: admin)

```http
PUT /api/v1/departments/1
Authorization: Bearer <token>
Content-Type: application/json

{
  "name": "IT Department",
  "description": "Updated description"
}
```

#### Xóa phòng ban (role: admin)

```http
DELETE /api/v1/departments/1
Authorization: Bearer <token>
```

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
// Định nghĩa interface
type UserRepository interface {
    Create(user *model.User) error
    FindByID(id uint) (*model.User, error)
    ...
}
// Implement interface
type userRepository struct { db *gorm.DB }
// Inject vào service
authSvc := service.NewAuthService(userRepo, roleRepo, &cfg.JWT)
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
| Language   | Go 1.21                 |
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
