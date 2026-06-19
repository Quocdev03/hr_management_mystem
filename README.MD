# Hệ Thống Quản Lý Nhân Sự (HR Management System)

![Go Version](https://img.shields.io/badge/Go-1.25-00ADD8?logo=go)
![Vue Version](https://img.shields.io/badge/Vue.js-3.x-4FC08D?logo=vuedotjs)
![Docker](https://img.shields.io/badge/Docker-Compose-2496ED?logo=docker)
![License](https://img.shields.io/badge/license-MIT-green)

Ứng dụng Quản Lý Nhân Sự (HRM) full-stack với Backend Go (Gin, GORM, MySQL, Redis) và Frontend Vue.js 3 (Vite, Pinia, Vue Router). Hỗ trợ quản lý nhân viên, phòng ban, tài khoản và phân quyền truy cập chi tiết.

---

## 📸 Giao diện ứng dụng

![Dashboard Preview](reviews/dashboard.png)

*(Xem thêm ảnh chụp màn hình các chức năng khác tại thư mục [reviews/](reviews/))*

---

## 🛠️ Tech Stack

- **Backend:** Go 1.25, Gin, GORM, Go-Redis, JWT
- **Frontend:** Vue.js 3 (Composition API), Vite, Pinia, Vue Router, Axios
- **Database:** MySQL 8.0
- **Cache / Rate Limiter:** Redis 7.2
- **DevOps:** Docker, Docker Compose

---

## 📖 Hướng Hướng Dẫn Cài Đặt & Chạy Ứng Dụng

Toàn bộ tài liệu hướng dẫn chi tiết đã được di chuyển sang thư mục [docs/](docs/):

- 🚀 **[Hướng dẫn chạy Local Development](docs/GETTING_STARTED.md)** — Bắt đầu nhanh với Docker và các file script `start-dev.bat` (Windows) hoặc qua Docker Compose.
- 🏭 **[Hướng dẫn deploy Production](docs/DEPLOYMENT_GUIDE.md)** — Cấu hình môi trường bảo mật, thiết lập SSL/TLS Nginx và deploy thực tế.

---

## 🏗️ Cấu Trúc Dự Án

```text
.
├── backend/                # Backend Go
│   ├── cmd/                # Entry points (server, migrate, seed)
│   ├── internal/           # Business logic
│   │   ├── config/         # Config (env, db, jwt, redis)
│   │   ├── handler/        # HTTP Handlers
│   │   ├── middleware/     # Auth, CORS, Cache, Rate limit
│   │   ├── model/          # GORM models
│   │   ├── repository/     # Data repository layer
│   │   ├── router/         # API routes
│   │   └── service/        # Service layer
│   ├── .air.toml           # Hot-reload config for Air
│   ├── Dockerfile          # Production Dockerfile
│   ├── Dockerfile.dev      # Development Dockerfile
│   ├── docker-entrypoint.sh
│   └── docker-entrypoint.dev.sh
├── frontend/               # Frontend Vue.js 3
│   ├── src/
│   │   ├── api/            # Axios instance
│   │   ├── components/     # Reusable components
│   │   ├── router/         # Routing
│   │   ├── store/          # Pinia stores
│   │   └── views/          # Views
│   ├── Dockerfile          # Nginx production build
│   └── vite.config.js
├── docs/                   # Tài liệu hướng dẫn
├── reviews/                # Ảnh chụp màn hình ứng dụng
├── scripts/                # Script tự động hóa
├── start-dev.bat           # Script chạy Development (CMD)
├── start-prod.bat          # Script chạy Production (CMD)
├── docker-compose.dev.yml  # Compose development
├── docker-compose.prod.yml # Compose production
└── database.sql            # Database schema backup
```

---

## 🔐 API Endpoints chính

API được versioned tại `/api/v1`. Hầu hết endpoints yêu cầu JWT Bearer token trong header.

| Method | Endpoint | Mô tả | Quyền truy cập |
|--------|----------|-------|----------------|
| `POST` | `/api/v1/auth/login` | Đăng nhập | Public |
| `GET` | `/api/v1/auth/me` | Lấy thông tin tài khoản hiện tại | User đăng nhập |
| `POST` | `/api/v1/auth/logout` | Đăng xuất (blacklist token) | User đăng nhập |
| `POST` | `/api/v1/auth/refresh` | Làm mới access token | User đăng nhập |
| `GET` | `/api/v1/employees` | Lấy danh sách nhân viên | `employee.read` |
| `POST` | `/api/v1/employees` | Tạo mới nhân viên | `employee.create` |
| `PUT/PATCH/DELETE` | `/api/v1/employees/:id` | Xem chi tiết, cập nhật, xóa nhân viên | `employee.read/update/delete` |
| `GET` | `/api/v1/departments` | Lấy danh sách phòng ban | `department.read` |
| `GET` | `/api/v1/users` | Lấy danh sách tài khoản | `user.read` (Admin) |
| `GET` | `/api/v1/dashboard/stats` | Lấy thống kê cho trang Dashboard | Đăng nhập |
| `GET` | `/api/v1/health` | Kiểm tra trạng thái của server | Public |

---

## 🗄️ Redis Caching & Rate Limiting Strategy

| Chức năng | Định dạng Key | TTL (Thời gian sống) | Mục đích |
|-----------|---------------|----------------------|----------|
| Token Blacklist | `blacklist:{token}` | Tương ứng thời hạn token | Đăng xuất an toàn |
| Rate Limit Login | `rate_limit:login:{IP}` | 1 phút | Chống brute-force đăng nhập |
| Refresh Token | `refresh_token:{userID}` | 7 ngày | Token rotation, duy trì session |
| Dashboard Cache | `dashboard:stats` | 1 giờ | Tăng tốc độ load trang Dashboard |
