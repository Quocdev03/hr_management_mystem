# 👨‍💻 Hướng Dẫn Chạy Project Lần Đầu

**Cho những người mới download source code về**

---

## 📋 Mục Tiêu

Hướng dẫn này dành cho ai mới tải source code và muốn chạy project trên máy tính cá nhân.

---

## 🎯 Bước 1: Chuẩn Bị (Cài Đặt Phần Mềm)

### Bước 1.1: Kiểm tra đã có những gì

Mở terminal/Command Prompt và chạy:

```bash
# Kiểm tra Docker
docker --version

# Kiểm tra Docker Compose
docker-compose --version

# Kiểm tra Git (để clone repo)
git --version
```

### Bước 1.2: Nếu chưa có, cài đặt:

**Windows:**

- Cài [Docker Desktop](https://www.docker.com/products/docker-desktop) (bao gồm Docker và Docker Compose)
- Cài [Git](https://git-scm.com/download/win)
- Cài [Notepad++](https://notepad-plus-plus.org/downloads/) hoặc editor yêu thích

**macOS:**

```bash
# Dùng Homebrew
brew install docker docker-compose git
brew install --cask docker
```

**Linux (Ubuntu/Debian):**

```bash
sudo apt update
sudo apt install docker.io docker-compose git
sudo usermod -aG docker $USER
```

---

## 📥 Bước 2: Tải Source Code

### Cách 1: Từ GitHub (nếu có repo)

```bash
# Tạo folder để chứa project
mkdir my-projects
cd my-projects

# Clone repository
git clone <URL_REPO>
cd HR_Management_System
```

### Cách 2: Tải file ZIP rồi giải nén

```bash
# Sau khi giải nén
cd HR_Management_System
```

---

## ⚙️ Bước 3: Setup Môi Trường

### Bước 3.1: Tạo file .env (QUAN TRỌNG!)

```bash
# Copy template
cp .env.example .env
```

### Bước 3.2: Chỉnh sửa file .env

**Trên Windows (Dùng Notepad++):**

```
Mở file .env bằng Notepad++
```

**Trên macOS/Linux:**

```bash
nano .env
# hoặc
code .env  # nếu dùng VS Code
```

### Bước 3.3: Điều chỉnh các giá trị quan trọng

Mở file `.env` và thay đổi những dòng này:

```env
# ===== PHẦN QUAN TRỌNG - THAY ĐỔI CÁC GIÁC TRỊ NÀY =====

# 1. Mật khẩu Database (tạo password ngẫu nhiên, min 16 ký tự)
DB_PASSWORD=your_strong_password_123456

# 2. Mật khẩu MySQL Root (tạo password ngẫu nhiên, min 16 ký tự)
MYSQL_ROOT_PASSWORD=root_strong_password_123456

# 3. JWT Secret (tạo string ngẫu nhiên, min 32 ký tự)
JWT_SECRET=your_jwt_secret_key_at_least_32_chars_long_12345

# 4. Các giá trị khác có thể giữ nguyên
APP_ENV=production
APP_SEED=false
APP_PORT=8080
DB_HOST=mysql
DB_PORT=3306
DB_USER=hrm_app_user
DB_NAME=hrm_db
REDIS_HOST=redis
REDIS_PORT=6379
```

**💡 Cách tạo password mạnh:**

```bash
# Tạo password 16 ký tự
openssl rand -base64 16

# Tạo JWT Secret 32 ký tự
openssl rand -base64 32
```

**Nếu không có OpenSSL:**

- Trên Windows: Dùng [Password Generator online](https://www.lastpass.com/features/password-generator)
- Hoặc tạo password: Kết hợp chữ hoa, chữ thường, số, ký tự đặc biệt

**Ví dụ:**

```env
DB_PASSWORD=MySecure@Pass2024#HRM
MYSQL_ROOT_PASSWORD=RootAdmin$2024#Secure
JWT_SECRET=aBcD1234EfGh5678IjKl9012MnOp3456QrSt7890
```

---

## 🔐 Bước 4: Tạo SSL Certificates (Chứng Chỉ HTTPS)

### Cho Development/Local (tự ký)

```bash
# Cấp quyền cho script
chmod +x scripts/generate-certs.sh

# Chạy script
./scripts/generate-certs.sh

# Sẽ tạo folder certs/ với 2 file:
# - certs/cert.pem
# - certs/key.pem
```

**Trên Windows (nếu chmod không hoạt động):**

```bash
# Trực tiếp chạy script bash
bash scripts/generate-certs.sh
```

---

## ✅ Bước 5: Kiểm Tra Cấu Hình (Tùy Chọn)

```bash
# Chạy validation script
chmod +x scripts/validate-production.sh
./scripts/validate-production.sh

# Nếu thấy "All checks passed!" là OK
```

**Windows:**

```bash
bash scripts/validate-production.sh
```

---

## 🚀 Bước 6: Chạy Project

### Chạy lần đầu (sẽ build Docker images)

```bash
# Chạy tất cả services (MySQL, Redis, Backend, Frontend)
docker-compose -f docker-compose.prod.yml up -d

# Chờ khoảng 30-60 giây cho services khởi động
sleep 30

# Kiểm tra trạng thái
docker-compose -f docker-compose.prod.yml ps

# Nên thấy: 4 containers đều "Up"
```

### Hoặc chạy với logs để xem quá trình

```bash
# Chạy và xem logs real-time
docker-compose -f docker-compose.prod.yml up

# Nhấn Ctrl+C để dừng
```

---

## ✨ Bước 7: Truy Cập Ứng Dụng

### Backend Health Check

```bash
# Kiểm tra API backend có chạy không
curl -k https://localhost/api/v1/health

# Nên thấy response:
# {"status":"ok","message":"HRM API đang chạy!"}
```

### Truy Cập Frontend

**Trên trình duyệt:**

- Đánh địa chỉ: `https://localhost`
- Sẽ thấy cảnh báo SSL (bình thường vì dùng self-signed certificate)
- Click "Advanced" → "Proceed to localhost"

### Đăng Nhập

```
Mở: https://localhost
Nhập username/password (xem trong database.sql)
Thường là:
  - Email: admin@example.com
  - Password: password123 (hoặc giá trị trong dữ liệu khởi tạo)
```

---

## 🛑 Bước 8: Dừng Project

```bash
# Dừng tất cả services
docker-compose -f docker-compose.prod.yml down

# Hoặc nếu muốn xóa toàn bộ dữ liệu (MySQL, Redis)
docker-compose -f docker-compose.prod.yml down -v
```

---

## 📊 Các Lệnh Hữu Ích

### Xem logs

```bash
# Xem tất cả logs
docker-compose -f docker-compose.prod.yml logs -f

# Xem logs của backend
docker-compose -f docker-compose.prod.yml logs -f backend

# Xem logs của frontend
docker-compose -f docker-compose.prod.yml logs -f frontend

# Xem logs của MySQL
docker-compose -f docker-compose.prod.yml logs -f mysql

# Nhấn Ctrl+C để thoát
```

### Kiểm tra trạng thái

```bash
# Xem trạng thái tất cả containers
docker-compose -f docker-compose.prod.yml ps

# Xem resource usage (CPU, Memory)
docker stats
```

### Restart service

```bash
# Restart tất cả
docker-compose -f docker-compose.prod.yml restart

# Restart một service cụ thể
docker-compose -f docker-compose.prod.yml restart backend
```

### Truy cập database

```bash
# Kết nối vào MySQL
docker-compose -f docker-compose.prod.yml exec mysql mysql -u hrm_app_user -p
# Nhập password từ DB_PASSWORD trong .env

# Hoặc với root:
docker-compose -f docker-compose.prod.yml exec mysql mysql -u root -p
# Nhập password từ MYSQL_ROOT_PASSWORD trong .env

# Sau khi kết nối:
USE hrm_db;
SHOW TABLES;
SELECT * FROM users;
EXIT;
```

### Truy cập Redis

```bash
# Kết nối vào Redis
docker-compose -f docker-compose.prod.yml exec redis redis-cli

# Các lệnh Redis:
PING          # Kiểm tra kết nối
INFO          # Thông tin Redis
DBSIZE        # Số lượng keys
KEYS *        # Liệt kê các keys
GET key_name  # Lấy giá trị
EXIT          # Thoát
```

---

## 🐛 Xử Lý Lỗi Thường Gặp

### Lỗi 1: Port đã được sử dụng

```
Error: bind: address already in use
```

**Giải pháp:**

```bash
# Tìm process đang dùng port
lsof -i :80      # Port 80
lsof -i :443     # Port 443
lsof -i :3306    # Port MySQL
lsof -i :6379    # Port Redis

# Kill process (hoặc dừng service khác)
kill -9 <PID>
```

**Windows:**

```bash
netstat -ano | findstr :80
taskkill /PID <PID> /F
```

### Lỗi 2: Docker daemon không chạy

```
Cannot connect to Docker daemon
```

**Giải pháp:**

- Mở Docker Desktop (trên Windows/macOS)
- Hoặc start Docker service (Linux):

```bash
sudo systemctl start docker
```

### Lỗi 3: SSL Certificate lỗi

```
SSL: CERTIFICATE_VERIFY_FAILED
```

**Giải pháp:**

```bash
# Tạo lại certificates
rm -rf certs/
./scripts/generate-certs.sh

# Hoặc sử dụng curl với -k flag
curl -k https://localhost/api/v1/health
```

### Lỗi 4: .env không được load

```bash
# Kiểm tra file .env tồn tại
ls -la .env

# Đảm bảo không có khoảng trắng quanh "="
# ✅ Đúng: DB_PASSWORD=mypass123
# ❌ Sai: DB_PASSWORD = mypass123
```

### Lỗi 5: Containers đang ở trạng thái "Exited"

```bash
# Xem logs để tìm nguyên nhân
docker-compose -f docker-compose.prod.yml logs

# Thường là lỗi:
# - .env không đúng
# - Port bị chiếm dụng
# - Disk space hết
```

---

## 📚 Tài Liệu Thêm

1. **[QUICK_START.md](QUICK_START.md)** - Hướng dẫn nhanh (5 phút)
2. **[DEPLOYMENT_GUIDE.md](DEPLOYMENT_GUIDE.md)** - Hướng dẫn chi tiết (production)
3. **[SECURITY_FIXES_SUMMARY.md](SECURITY_FIXES_SUMMARY.md)** - Các cải thiện bảo mật
4. **[README.md](README.md)** - Tổng quan project

---

## ✅ Checklist Hoàn Thành

Sau khi hoàn thành, bạn nên có:

- [x] Docker và Docker Compose cài đặt
- [x] Source code tải về
- [x] File .env tạo từ .env.example
- [x] Passwords mạnh được thiết lập
- [x] SSL certificates được tạo
- [x] Project chạy lên (4 containers)
- [x] Có thể truy cập https://localhost
- [x] Backend health check thành công
- [x] Có thể đăng nhập vào ứng dụng

---

## 🆘 Cần Giúp?

**Nếu gặp vấn đề:**

1. Kiểm tra lại các file:
   - `.env` có tồn tại?
   - `certs/cert.pem` và `certs/key.pem` có tồn tại?
   - File có quyền đọc/ghi không?

2. Xem logs:

   ```bash
   docker-compose -f docker-compose.prod.yml logs
   ```

3. Xem trạng thái containers:

   ```bash
   docker-compose -f docker-compose.prod.yml ps
   ```

4. Restart lại từ đầu:
   ```bash
   docker-compose -f docker-compose.prod.yml down -v
   # (Sửa lỗi)
   docker-compose -f docker-compose.prod.yml up -d
   ```

---

## 🎉 Xong!

Nếu bạn thấy ứng dụng chạy được, chúc mừng! 🎊

Bây giờ bạn có thể:

- Thử nghiệm các tính năng
- Phát triển thêm
- Deploy lên production (xem DEPLOYMENT_GUIDE.md)

---

**Phiên bản:** 1.0  
**Cập nhật lần cuối:** 2026-06-14  
**Trạng thái:** ✅ Hoàn chỉnh
