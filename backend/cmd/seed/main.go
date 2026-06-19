package main

import (
	"chiquoc_hocgolang/internal/config"
	"chiquoc_hocgolang/internal/utils"
	"os"
)

func main() {
	utils.Info("[SEED] Loading environment...")
	cfg := config.Load()

	utils.Info("[SEED] Connecting to database...")
	db := config.InitDB(&cfg.Database)
	utils.Info("[SEED] Database connected.")

	// Chạy seed — idempotent: chạy nhiều lần không tạo dữ liệu trùng
	utils.Info("[SEED] Running seed data...")
	if err := config.SeedDataWithError(db); err != nil {
		utils.Error("[SEED] Seed failed: %v", err)
		os.Exit(1)
	}

	utils.Info("[SUCCESS] Seed completed.")
}
