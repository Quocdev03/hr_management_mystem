package main

import (
	"chiquoc_hocgolang/internal/config"
	"log"
)

func main() {
	// Load config từ .env
	cfg := config.Load()

	// Bước 1: Tạo database nếu chưa có
	log.Println("=== Bước 1: Tạo database ===")
	config.CreateDatabase(&cfg.Database)

	// Bước 2: Kết nối database
	log.Println("=== Bước 2: Kết nối database ===")
	db := config.InitiDB(&cfg.Database)

	// Bước 3: Chạy migrations (tạo / cập nhật bảng)
	log.Println("=== Bước 3: Chạy migrations ===")
	config.RunMigrations(db)

	// Bước 4: Chèn dữ liệu mẫu
	log.Println("=== Bước 4: Seed dữ liệu mẫu ===")
	config.SeedData(db)

	log.Println("=== Setup hoàn tất! ===")
}
