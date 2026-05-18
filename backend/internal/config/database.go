/*
Mục tiêu file này
- Kết nối database
- Cấu hình connection pool
- Auto migration
*/
package config

import (
	"chiquoc_hocgolang/internal/model"
	"fmt"
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

	dsnWithoutDB := fmt.Sprintf("%s:%s@tcp(%s:%s)/?charset=utf8mb4&parseTime=True&loc=Local", cfg.User, cfg.Password, cfg.Host, cfg.Port)
	tempDB, err := gorm.Open(mysql.Open(dsnWithoutDB), &gorm.Config{
		Logger: glogger.Default.LogMode(glogger.Silent),
	})
	if err != nil {
		log.Fatalf("Không thể kết nối MySQL server: %v", err)
	}

	// createDBQuery := fmt.Sprintf("CREATE DATABASE IF NOT EXISTS `%s` CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci", cfg.DBName)
	// if err := tempDB.Exec(createDBQuery).Error; err != nil {
	// 	log.Fatalf("Không thể tạo database: %v", err)
	// }

	if sqlTempDB, err := tempDB.DB(); err == nil {
		sqlTempDB.Close()
	}

	db, err := gorm.Open(mysql.Open(cfg.DSN()), &gorm.Config{
		Logger: glogger.Default.LogMode(logLevel),
	})
	if err != nil {
		log.Fatalf("Kết nối cơ sở dữ liệu bị lỗi: %v", err)
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
	runMigrations(db)

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
