package main

import (
	"chiquoc_hocgolang/internal/config"
	"chiquoc_hocgolang/internal/utils"
	"os"
)

func main() {
	utils.Info("[MIGRATE] Loading environment...")
	cfg := config.Load()

	utils.Info("[MIGRATE] Connecting to database...")
	utils.Info("[MIGRATE] DB Host: %s:%s / DB: %s", cfg.Database.Host, cfg.Database.Port, cfg.Database.DBName)

	// Bước 1: Tạo database nếu chưa có
	config.CreateDatabase(&cfg.Database)
	utils.Info("[MIGRATE] Database ensured.")

	// Bước 2: Kết nối database
	db := config.InitDB(&cfg.Database)
	utils.Info("[MIGRATE] Database connected.")

	// Bước 3: Chạy migrations
	utils.Info("[MIGRATE] Running migrations...")
	if err := config.RunMigrationsWithError(db); err != nil {
		utils.Error("[MIGRATE] Migration failed: %v", err)
		os.Exit(1)
	}

	utils.Info("[SUCCESS] Migration completed.")
}
