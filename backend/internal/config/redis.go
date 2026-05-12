package config

import (
	"context"
	"log"

	"github.com/redis/go-redis/v9"
)

// InitRedis khởi tạo kết nối Redis
func InitRedis(cfg *RedisConfig) *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr:     cfg.RedisAddr(),
		Password: cfg.Password,
		DB:       cfg.DB,
	})

	// Ping để kiểm tra kết nối
	ctx := context.Background()
	if _, err := client.Ping(ctx).Result(); err != nil {
		// Không fatal vì Redis là optional (cache layer)
		log.Printf("Lỗi kết nối tới Redis: %v - đang chạy mà không có cache!", err)
		return nil
	}

	log.Println("Kết nối với Redis hoàn tất!")
	return client

}
