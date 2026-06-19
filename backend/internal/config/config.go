/*
Mục tiêu file này:
- Quản lý toàn bộ cấu hình ứng dụng
- Đọc biến môi trường (.env)
- Cung cấp config cho Database, JWT, App
*/
package config

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	App      AppConfig
	Database DatabaseConfig
	JWT      JWTConfig
	Redis    RedisConfig
}

type AppConfig struct {
	Port string
	Env  string
	Seed bool
}

type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	Env      string
}

type JWTConfig struct {
	SecretKey        string
	ExpireHour       int
	RefreshExpireDay int
}

type RedisConfig struct {
	Host     string
	Port     string
	Password string
	DB       int
}

func Load() *Config {
	// Chỉ load .env khi đang chạy local
	if err := godotenv.Load(); err != nil {
		log.Println("không tìm thấy file .env, đang dùng biến môi trường")
	}

	env := getEnv("APP_ENV", "production")

	// Chuyển string -> int vì env chỉ chứa string
	expireHour, _ := strconv.Atoi(getEnv("JWT_EXPIRE_HOUR", "24"))
	refreshExpireDay, _ := strconv.Atoi(getEnv("JWT_REFRESH_EXPIRE_DAY", "7"))
	redisDB, _ := strconv.Atoi(getEnv("REDIS_DB", "0"))

	return &Config{
		App: AppConfig{
			Port: getEnv("APP_PORT", "8080"),
			Env:  env,
			Seed: getEnv("APP_SEED", "false") == "true",
		},
		Database: DatabaseConfig{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnv("DB_PORT", "3306"),
			User:     getEnv("DB_USER", "root"),
			Password: requireEnv("DB_PASSWORD", "password123", env),
			DBName:   getEnv("DB_NAME", "hrm_db"),
			Env:      env,
		},
		JWT: JWTConfig{
			SecretKey:        requireEnv("JWT_SECRET", "your-super-secret-key-min-32-chars-change-in-production", env),
			ExpireHour:       expireHour,
			RefreshExpireDay: refreshExpireDay,
		},
		Redis: RedisConfig{
			Host:     getEnv("REDIS_HOST", "localhost"),
			Port:     getEnv("REDIS_PORT", "6379"),
			Password: getEnv("REDIS_PASSWORD", ""),
			DB:       redisDB,
		},
	}
}

// DSN trả vè chuỗi kết nối MySQL theo format GORM
func (d *DatabaseConfig) DSN() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		d.User, d.Password, d.Host, d.Port, d.DBName)
}

// Hàm đọc biến môi trường, trả ra value mặc định nếu không tìm thấy
func getEnv(key, defaultVal string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return defaultVal
}

// requireEnv đọc biến môi trường bắt buộc.
// Trong production nếu thiếu → log.Fatal ngay lập tức, không fallback.
// Trong development/test thì dùng giá trị mặc định.
func requireEnv(key, devDefault, env string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	if env == "production" {
		log.Fatalf("[FATAL] Biến môi trường bắt buộc '%s' không được đặt trong production", key)
	}
	return devDefault
}
