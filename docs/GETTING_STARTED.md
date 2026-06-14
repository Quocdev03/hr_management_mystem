👋 **CHÀO MỪNG!** Bạn vừa tải source code project HR Management System

---

# 🚀 GETTING STARTED

**Hướng dẫn nhanh và chi tiết cho người mới**

---

## ⚡ Chạy Ngay (3 Phút)

Nếu bạn chỉ muốn chạy project ngay:

```bash
# 1. Setup environment
cp .env.example .env
# Edit 3 dòng quan trọng trong .env:
#  - DB_PASSWORD=strong_pass_16_chars
#  - MYSQL_ROOT_PASSWORD=root_pass_16_chars
#  - JWT_SECRET=jwt_secret_32_chars

# 2. Tạo SSL certificates
chmod +x scripts/generate-certs.sh
./scripts/generate-certs.sh

# 3. Chạy project
docker-compose -f docker-compose.prod.yml up -d

# 4. Chờ ~30 giây
sleep 30

# 5. Truy cập
# Frontend: https://localhost
# Backend Health: curl -k https://localhost/api/v1/health
```

✅ **Xong!** Project chạy được.

---

## 📚 Chọn Hướng Dẫn Phù Hợp

### 🏃 "Tôi chỉ muốn chạy project ngay"

⏱️ **Mất 3 phút**

```bash
cp .env.example .env
# Edit .env: change DB_PASSWORD, MYSQL_ROOT_PASSWORD, JWT_SECRET

chmod +x scripts/generate-certs.sh
./scripts/generate-certs.sh

docker-compose -f docker-compose.prod.yml up -d
sleep 30

# Truy cập: https://localhost
# Backend: curl -k https://localhost/api/v1/health
```

✅ **Xong!** Project chạy được.

---

### 🧑‍🎓 "Tôi muốn hiểu từng bước"

⏱️ **Mất 20 phút**

👉 **[Đọc: SETUP_FOR_NEWCOMERS.md](SETUP_FOR_NEWCOMERS.md)**

Bao gồm:

- Cài đặt phần mềm yêu cầu
- Setup environment từng bước
- Giải thích từng lệnh
- Xử lý lỗi thường gặp
- Lệnh hữu ích

---

### 🎯 "Tôi muốn biết project này làm gì"

⏱️ **Mất 30 phút**

👉 **[Đọc: README.md](README.md)**

---

### 🔐 "Tôi muốn biết về bảo mật"

⏱️ **Mất 5 phút**

👉 **Xem phần "🗄️ Chiến Lược Redis"** trong [README.md](../README.MD)

Bao gồm:
- Fix `ParseUnverified` vulnerability (JWT)
- Caching strategy (Dashboard cache 1h, loại bỏ GET cache 95% miss rate)
- Rate limiting và Token Blacklist

---

### 🌍 "Tôi muốn deploy lên production"

⏱️ **Mất 60 phút** (bắt buộc đọc)

👉 **[Đọc: DEPLOYMENT_GUIDE.md](DEPLOYMENT_GUIDE.md)**

---

## 🎯 Lộ Trình Nhanh

```
Lần 1: Chạy ngay (3 phút)
├─ GETTING_STARTED.md (Quick section)
└─ ✅ Chạy được

Lần 2: Hiểu chi tiết (20 phút)
├─ SETUP_FOR_NEWCOMERS.md
└─ ✅ Hiểu rõ

Lần 3: Nâng cao (15 phút)
├─ README.md (Tech stack + API + Redis strategy)
└─ ✅ Hiểu sâu

Lần 4: Production (60 phút)
├─ DEPLOYMENT_GUIDE.md
└─ ✅ Deploy thành công
```

---

## 📋 Tài Liệu Có Sẵn

| Tài Liệu                                                   | Mục Đích                 | Thời Gian |
| ---------------------------------------------------------- | ------------------------ | --------- |
| **[SETUP_FOR_NEWCOMERS.md](SETUP_FOR_NEWCOMERS.md)**       | Setup chi tiết từng bước | 20 min    |
| **[DEPLOYMENT_GUIDE.md](DEPLOYMENT_GUIDE.md)**             | Deploy lên production    | 60 min    |
| **[README.md](../README.MD)**                              | Tổng quan + Redis strategy | 10 min  |

---

## ✅ Yêu Cầu

- Docker & Docker Compose
- Git (tùy chọn)
- Editor (VS Code, Notepad++, etc.)

```bash
docker --version
docker-compose --version
```

---

## 📞 Gặp Vấn Đề?

| Lỗi                 | Giải Pháp                                        |
| ------------------- | ------------------------------------------------ |
| Port đang dùng      | `docker-compose -f docker-compose.prod.yml down` |
| SSL certificate     | `./scripts/generate-certs.sh`                    |
| .env không tìm      | `cp .env.example .env`                           |
| Containers "Exited" | `docker-compose -f docker-compose.prod.yml logs` |

👉 **Chi tiết xem:** [SETUP_FOR_NEWCOMERS.md](SETUP_FOR_NEWCOMERS.md)

---

## 🎉 Bắt Đầu Thôi!

1. **Chỉ có 3 phút?** → Chạy lệnh quick ở trên
2. **Có 20 phút?** → Đọc [SETUP_FOR_NEWCOMERS.md](SETUP_FOR_NEWCOMERS.md)
3. **Muốn deploy?** → Đọc [DEPLOYMENT_GUIDE.md](DEPLOYMENT_GUIDE.md)

---

**Cập nhật:** 2026-06-14 | **Phiên bản:** 1.1 | **Trạng thái:** ✅
