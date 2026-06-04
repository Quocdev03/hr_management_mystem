#!/bin/sh

# Mặc định nếu không được cấu hình qua env
DB_HOST=${DB_HOST:-mysql}
DB_PORT=${DB_PORT:-3306}
REDIS_HOST=${REDIS_HOST:-redis}
REDIS_PORT=${REDIS_PORT:-6379}

echo "================================================================="
echo "⚙️ Khởi tạo môi trường Backend..."
echo "================================================================="

# Chờ MySQL sẵn sàng
echo "▶ Đang chờ MySQL kết nối tại $DB_HOST:$DB_PORT..."
while ! nc -z $DB_HOST $DB_PORT; do
  sleep 1
done
echo "✔ MySQL đã hoạt động!"

# Chờ Redis sẵn sàng
echo "▶ Đang chờ Redis kết nối tại $REDIS_HOST:$REDIS_PORT..."
while ! nc -z $REDIS_HOST $REDIS_PORT; do
  sleep 1
done
echo "✔ Redis đã hoạt động!"

# Thực hiện setup cơ sở dữ liệu (migration và seed)
echo "▶ Đang kiểm tra và chạy setup DB (migrations & seed)..."
./setup

# Khởi chạy server
echo "🚀 Khởi động API Server..."
echo "================================================================="
exec ./server
