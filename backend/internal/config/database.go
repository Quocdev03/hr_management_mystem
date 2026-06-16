/*
Mục tiêu file này
- Kết nối database
- Cấu hình connection pool
- Auto migration
*/
package config

import (
	"log"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	glogger "gorm.io/gorm/logger"
)

func InitDB(cfg *DatabaseConfig) *gorm.DB {
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
		log.Fatalf("kết nối cơ sở dữ liệu bị lỗi: %v", err)
	}

	// Cấu hình connection pool
	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalf("không thể lấy sql.DB: %v", err)
	}

	// Tối đa 10 kết nối cùng lúc, giữ 5 kết nối, không hết hạn
	sqlDB.SetMaxOpenConns(10)
	sqlDB.SetMaxIdleConns(5)
	sqlDB.SetConnMaxLifetime(1)

	log.Println("kết nối database và tạo các bảng thành công")
	return db
}
