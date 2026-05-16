/*
Mục tiêu file này
- Kết nối database
- Cấu hình connection pool
- Auto migration
*/
package config

import (
	"chiquoc_hocgolang/internal/model"
	"log"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	glogger "gorm.io/gorm/logger"
)

func InitiDB(cfg *DatabaseConfig) *gorm.DB {
	// 1. Set log level theo môi trường cho GORM
	var logLevel glogger.LogLevel

	if cfg.Env == "development" {
		logLevel = glogger.Info
	} else {
		logLevel = glogger.Warn
	}

	db, err := gorm.Open(mysql.Open(cfg.DSN()), &gorm.Config{
		Logger: glogger.Default.LogMode(logLevel),
	})
	if err != nil {
		log.Fatal("Kết nối cơ sở dữ liệu bị lỗi!")
	}

	// Cấu hình connection pool
	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalf("Không thể lấy sql.DB: %v", err)
	}

	// Tối đa 10 kết nối cùng lúc, giữ 5 kết nối, không hết hạn
	sqlDB.SetMaxOpenConns(10)
	sqlDB.SetMaxIdleConns(5)
	sqlDB.SetConnMaxLifetime(0)

	// Auto migrate: tự động tạo/cập nhật bảng dựa trên struct
	// runMigrations(db)

	log.Println("Kết nối database và tạo các bảng thành công!")
	return db
}
func runMigrations(db *gorm.DB) {
	err := db.AutoMigrate(
		&model.Role{},
		&model.Permission{},
		&model.User{},
		&model.Department{},
		&model.Employee{},
	)
	if err != nil {
		log.Fatalf("Tạo bảng thất bại: %v", err)
	}
}
