# 🚀 Getting Started — HR Management System

> **Tài liệu hướng dẫn cài đặt và khởi chạy dự án trong môi trường phát triển (Development).**

---

## ✅ Yêu cầu hệ thống

Chỉ cần cài đặt **Docker Desktop** trên máy của bạn:

- **Windows / macOS:** Tải và cài đặt tại [Docker Desktop](https://www.docker.com/products/docker-desktop)
- **Linux:** Cài đặt thông qua package manager: `sudo apt install docker.io docker-compose-plugin`

---

## ⚡ Khởi chạy dự án (Development)

### Trên Windows (Sử dụng Command Prompt / CMD)

Mở CMD trong thư mục gốc của dự án và chạy:

```cmd
:: Lần đầu tiên chạy (để thiết lập DB, chạy migration và seed dữ liệu mẫu)
start-dev.bat setup

:: Các lần tiếp theo chạy (chỉ khởi động hệ thống)
start-dev.bat
```

### Trên macOS / Linux (Sử dụng Terminal)

Mở Terminal trong thư mục gốc của dự án và chạy:

```bash
# Lần đầu tiên chạy: Khởi chạy docker compose, sau đó chạy migrate và seed
docker compose -f docker-compose.dev.yml up -d
docker compose -f docker-compose.dev.yml run --rm migrate
docker compose -f docker-compose.dev.yml run --rm seed

# Các lần tiếp theo chạy
docker compose -f docker-compose.dev.yml up -d
```

---

## 🌍 Các cổng truy cập (Ports)

Sau khi hệ thống khởi chạy thành công:

- **Frontend (Vue 3 + Vite HMR):** [http://localhost:5173](http://localhost:5173)
- **Backend API (Go + Gin):** [http://localhost:8080/api/v1/health](http://localhost:8080/api/v1/health)

---

## 🔑 Tài khoản đăng nhập demo

Hệ thống đã được tự động seed các tài khoản mẫu sau trong môi trường dev:

| Vai trò (Role) | Email đăng nhập | Mật khẩu (Password) |
|----------------|----------------|---------------------|
| **Admin** | `chiquoc23AD@company.vn` | `chiquoc23AD` |
| **HR (Nhân sự)** | `chiquoc23HR@company.vn` | `chiquoc23HR` |
| **Employee (Nhân viên)** | `chiquoc23EMP@company.vn` | `chiquoc23EMP` |

---

## 📋 Các lệnh quản lý khác (Windows CMD)

| Mục đích | Lệnh chạy |
|----------|-----------|
| **Dừng hệ thống** | `start-dev.bat down` |
| **Xem logs** | `start-dev.bat logs` |
| **Xem trạng thái container** | `start-dev.bat status` |

---

## 🔧 Cấu hình biến môi trường (.env)

Khi chạy `start-dev.bat`, file `.env` sẽ tự động được copy từ `.env.example` nếu chưa tồn tại.
Để tùy chỉnh cổng kết nối hoặc thông số DB, bạn có thể chỉnh sửa trực tiếp trong file `.env`:

```env
APP_PORT=8080        # Cổng chạy Backend API
VITE_PORT=5173       # Cổng chạy Frontend
DB_PASSWORD=change_me_strong_db_password   # Mật khẩu MySQL
JWT_SECRET=change_me_min_32_chars_cryptographically_random_key # Khóa ký JWT
```

---

## 🐛 Xử lý một số lỗi thường gặp

| Hiện tượng | Nguyên nhân & Cách khắc phục |
|------------|------------------------------|
| `Docker is not running` | Docker Desktop chưa được mở hoặc đang khởi động. Hãy mở Docker Desktop và đợi nó báo xanh. |
| Cổng `5173` hoặc `8080` bị chiếm | Một ứng dụng khác đang dùng cổng này. Bạn có thể sửa `VITE_PORT` / `APP_PORT` tương ứng trong file `.env` rồi khởi động lại. |
| Container dừng bất thường | Chạy lệnh `start-dev.bat logs` để xem thông tin lỗi của container đó. |
