# Hệ Thống Quản Lý Nhân Sự (HR Management System)

![Go Version](https://img.shields.io/badge/Go-1.26.2+-00ADD8?logo=go)
![Vue Version](https://img.shields.io/badge/Vue.js-3.x-4FC08D?logo=vuedotjs)
![License](https://img.shields.io/badge/license-MIT-green)

Một ứng dụng Quản Lý Nhân Sự (HRM) toàn diện (full-stack) được xây dựng với Backend bằng Go và Frontend bằng Vue.js 3. Hệ thống này cung cấp các chức năng cốt lõi để quản lý nhân viên, phòng ban và tài khoản người dùng cùng với phân quyền truy cập.

> 📖 **Hướng dẫn Setup:** Xem [GETTING_STARTED.md](/docs/GETTING_STARTED.md)

---

## 🛠️ Công Nghệ Sử Dụng (Tech Stack)

- **Backend:** Go, Gin, GORM, Go-Redis, JWT
- **Frontend:** Vue.js 3 (Composition API), Vite, Pinia, Vue Router, Axios
- **Database:** MySQL
- **Cache/Rate Limiter:** Redis
- **DevOps:** Docker, Docker Compose

## ✨ Tính Năng

### Backend (Go)

- **RESTful API:** Xây dựng với Gin, GORM ORM cho tương tác MySQL hiệu quả.
- **Xác Thực JWT:** Refresh Token rotation + Token Blacklist lưu Redis với TTL tự động. Logout an toàn — chỉ dùng verified claims, không dùng `ParseUnverified`.
- **Phân Quyền (RBAC):** Role-based access control cho `admin`, `hr`, `employee` roles.
- **Caching thông minh:** Dashboard stats cache TTL 1h (hit rate ~99%). GET endpoints truy vấn DB trực tiếp (index tốt, <100ms). Tránh URL-based cache với 95% miss rate.
- **Rate Limiting:** Redis rate limit login — 5 lượt/phút/IP (Fixed Window Counter).
- **Docker & Environment Config:** Multi-environment support với `.env`, dịch vụ containerize hoàn toàn.

### Frontend (Vue.js)

- **Vue 3 + Composition API:** Framework hiện đại với Vite build tool cho development cực nhanh.
- **State Management & Routing:** Pinia quản lý trạng thái + Vue Router với role-based route guards.
- **Giao Diện CRUD & Dashboard:** Quản lý Nhân Viên, Phòng Ban, Người Dùng + bảng thống kê tập trung.
- **Token Management:** Axios interceptor + Silent Refresh (tự động retry 401) + Responsive UI đa kích thước.
- **Reusable Components:** Modal, Dialog, Skeleton loader, UI components để tối ưu UX.

## 📸 Ảnh Chụp Màn Hình (Screenshots)

|                Trang Đăng Nhập                |            Bảng Điều Khiển (Dashboard)            |
| :-------------------------------------------: | :-----------------------------------------------: |
| ![Trang Đăng Nhập](./images_review/login.png) | ![Bảng Điều Khiển](./images_review/dashboard.png) |

|                 Quản Lý Nhân Viên                  |                  Quản Lý Phòng Ban                   |
| :------------------------------------------------: | :--------------------------------------------------: |
| ![Quản Lý Nhân Viên](./images_review/employee.png) | ![Quản Lý Phòng Ban](./images_review/department.png) |

|               Quản Lý Người Dùng                |                 Hồ Sơ Cá Nhân                 |
| :---------------------------------------------: | :-------------------------------------------: |
| ![Quản Lý Người Dùng](./images_review/user.png) | ![Hồ Sơ Cá Nhân](./images_review/profile.png) |

## ⚙️ Cấu Trúc Dự Án (Project Structure)

Kho mã nguồn được sắp xếp theo cấu trúc monorepo bao gồm `backend` và `frontend`.

```text
.
├── backend/                # Ứng dụng Backend Go
│   ├── cmd/                # Các điểm khởi chạy chính (server, setup)
│   ├── internal/           # Logic lõi của ứng dụng
│   │   ├── config/         # Cấu hình (env, db, jwt)
│   │   ├── handler/        # Các HTTP Handler (controllers)
│   │   ├── middleware/     # API Middleware (auth, cors, cache)
│   │   ├── model/          # Các Model và DTO
│   │   ├── repository/     # Tầng truy xuất dữ liệu (GORM)
│   │   ├── router/         # Khai báo các API route
│   │   └── service/        # Logic nghiệp vụ (Business logic)
│   ├── go.mod
│   └── ...
├── frontend/               # Ứng dụng Frontend Vue.js
│   ├── public/
│   ├── src/
│   │   ├── api/            # Cấu hình và Axios instance
│   │   ├── assets/         # Tài nguyên tĩnh (CSS, hình ảnh)
│   │   ├── components/     # Component Vue dùng lại
│   │   ├── helpers/        # Các hàm hỗ trợ và composable
│   │   ├── layout/         # Giao diện khung (layout) chính
│   │   ├── router/         # Cấu hình Vue Router
│   │   ├── store/          # Quản lý trạng thái bằng Pinia
│   │   └── views/          # Các component Trang (Pages)
│   ├── package.json
│   └── ...
├── docker-compose.yml      # Cấu hình dịch vụ Docker (MySQL, Redis)
├── Makefile                # Các lệnh hỗ trợ quá trình phát triển
└── database.sql            # File định nghĩa Schema SQL thuần
```

## 🔐 Danh Sách API (API Endpoints)

API được đánh dấu phiên bản `/api/v1`. Hầu hết các endpoint yêu cầu token xác thực JWT qua tiêu đề `Bearer`.

- **Xác Thực (Authentication - `/auth`)**
   - `POST /login`: Đăng nhập và nhận JWT.
   - `GET /me`: Lấy thông tin hồ sơ của người dùng đang đăng nhập.
   - `POST /logout`: Đăng xuất, blacklist access token (Redis TTL), xoá refresh token.
   - `POST /refresh`: Lấy access token mới bằng refresh token (token rotation).

- **Nhân Viên (Employees - `/employees`)**
   - `GET /`: Lấy danh sách nhân viên có phân trang.
   - `POST /`: Tạo nhân viên mới (admin, hr).
   - `GET /:id`: Lấy thông tin một nhân viên theo ID.
   - `PUT /:id`: Cập nhật nhân viên đầy đủ (admin, hr).
   - `PATCH /:id`: Cập nhật nhân viên một phần (admin, hr).
   - `DELETE /:id`: Xóa nhân viên (chỉ admin).

- **Phòng Ban (Departments - `/departments`)**
   - `GET /`: Lấy danh sách tất cả phòng ban.
   - `POST /`: Tạo phòng ban mới (chỉ admin).
   - `GET /:id`: Lấy thông tin một phòng ban theo ID.
   - `PUT /:id`: Cập nhật phòng ban đầy đủ (chỉ admin).
   - `PATCH /:id`: Cập nhật phòng ban một phần (chỉ admin).
   - `DELETE /:id`: Xóa phòng ban (chỉ admin).

- **Người Dùng (Users - `/users`)**
   - `GET /`: Lấy danh sách người dùng có phân trang (chỉ admin).
   - `POST /`: Tạo người dùng mới (chỉ admin).
   - `GET /available`: Lấy các user chưa được liên kết với hồ sơ nhân viên (chỉ admin).
   - `GET /:id`: Lấy thông tin một người dùng theo ID (chỉ admin).
   - `PUT /:id`: Cập nhật người dùng đầy đủ (chỉ admin).
   - `PATCH /:id`: Cập nhật người dùng một phần (chỉ admin).
   - `DELETE /:id`: Xóa người dùng (chỉ admin).

- **Bảng Điều Khiển (Dashboard - `/dashboard`)**
   - `GET /stats`: Lấy thống kê hệ thống (số lượng nhân viên, phòng ban, v.v.). **Cache Redis 1h**, tự động invalidate khi dữ liệu thay đổi.

## 🗄️ Chiến Lược Redis

| Tính Năng             | Key Format                    | TTL       | Ghi Chú                           |
| --------------------- | ----------------------------- | --------- | --------------------------------- |
| Token Blacklist       | `blacklist:{token}`           | Tự động   | Bắt buộc — bảo mật logout        |
| Rate Limit Login      | `rate_limit:login:{IP}`       | 1 phút    | Bắt buộc — chống brute-force     |
| Refresh Token Storage | `refresh_token:{userID}`      | 7 ngày    | Bắt buộc — token rotation        |
| Dashboard Stats Cache | `dashboard:stats`             | 1 giờ     | Tốt — hit rate ~99%, key cố định |
| ~~GET Response Cache~~| ~~`cache:{URL+params}`~~      | ~~Đã xóa~~| Xấu — 95% cache miss rate        |

## 🚀 Bắt Đầu

👉 **[Hướng dẫn setup chi tiết tại: docs/GETTING_STARTED.md](docs/GETTING_STARTED.md)**

**Quick 3-step (3 phút):**

```bash
cp .env.example .env          # Setup environment
chmod +x scripts/generate-certs.sh && ./scripts/generate-certs.sh  # SSL
docker-compose -f docker-compose.prod.yml up -d  # Run
```

**Xem thêm:**

- 📖 [docs/SETUP_FOR_NEWCOMERS.md](docs/SETUP_FOR_NEWCOMERS.md) - Chi tiết từng bước
- 🚀 [docs/DEPLOYMENT_GUIDE.md](docs/DEPLOYMENT_GUIDE.md) - Deploy production
