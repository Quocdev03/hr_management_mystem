package config

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

// Config chứa toàn bộ cấu hình ứng dụng
type Config struct {
	App      AppConfig
	Database DatabaseConfig
	JWT      JWTConfig
	Redis    RedisConfig
}

type AppConfig struct {
	Port string
	Env  string
}

type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
}

type JWTConfig struct {
	SecretKey  string
	ExpireHour int
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
		log.Println("Không tìm thấy file .env, đang dùng biến môi trường!")
	}
	// Chuyển string -> int vì env chỉ chứa string
	expireHour, _ := strconv.Atoi(getEnv("JWT_EXPIRE_HOUR", "24"))
	redisDB, _ := strconv.Atoi(getEnv("REDIS_DB", "0"))

	return &Config{
		App: AppConfig{
			Port: getEnv("APP_PORT", "8080"),
			Env:  getEnv("APP_ENV", "development"),
		},
		Database: DatabaseConfig{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnv("DB_PORT", "3306"),
			User:     getEnv("DB_USER", "root"),
			Password: getEnv("DB_PASSWORD", "password123"),
			DBName:   getEnv("DB_NAME", "hrm_db"),
		},
		JWT: JWTConfig{
			SecretKey:  getEnv("JWT_SECRET", "your-super-secret-key-min-32-chars-change-in-production"),
			ExpireHour: expireHour,
		},
		Redis: RedisConfig{
			Host:     getEnv("REDIS_HOST", "localhost"),
			Port:     getEnv("REDIS_PORT", "6379"),
			Password: getEnv("REDIS_PASSWORD", "password"),
			DB:       redisDB,
		},
	}
}

// DSN trả vè chuỗi kết nối MySQL theo format GORM
func (d *DatabaseConfig) DSN() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		d.User, d.Password, d.Host, d.Port, d.DBName)
}

// RedisAddr trả về địa chỉ Redis
func (r *RedisConfig) RedisAddr() string {
	return fmt.Sprintf("%s:%s", r.Host, r.Port)
}

// Hàm đọc biến môi trường, trả ra value mặc định nếu không tìm thấy
func getEnv(key, defaultVal string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return defaultVal
}
