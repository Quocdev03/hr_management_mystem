/*
Mục tiêu file này
- Kết nối database
- Cấu hình connection pool
- Auto migration
*/
package config

import (
	"log"
	"time"

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

	var db *gorm.DB
	var err error
	maxRetries := 15
	for i := 1; i <= maxRetries; i++ {
		db, err = gorm.Open(mysql.Open(cfg.DSN()), &gorm.Config{
			Logger: glogger.Default.LogMode(logLevel),
		})
		if err == nil {
			break
		}
		log.Printf("đang chờ kết nối database... (Lần thử %d/%d): %v\n", i, maxRetries, err)
		time.Sleep(2 * time.Second)
	}
	if err != nil {
		log.Fatalf("kết nối cơ sở dữ liệu bị lỗi: %v", err)
	}

	// Cấu hình connection pool
	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalf("không thể lấy sql.DB: %v", err)
	}

	// Tối đa 10 kết nối cùng lúc, giữ 5 kết nối, mỗi connection sống tối đa 1 giờ
	sqlDB.SetMaxOpenConns(10)
	sqlDB.SetMaxIdleConns(5)
	sqlDB.SetConnMaxLifetime(time.Hour)

	log.Println("kết nối database và tạo các bảng thành công")
	return db
}
