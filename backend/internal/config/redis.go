package config

import (
	"chiquoc_hocgolang/internal/utils"
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

// InitRedis khởi tạo Redis Client và ping thử kết nối
func InitRedis(cfg *RedisConfig) *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", cfg.Host, cfg.Port),
		Password: cfg.Password,
		DB:       cfg.DB,
	})

	// Test connection khi khởi chạy với timeout 5s
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := rdb.Ping(ctx).Result()
	if err != nil {
		utils.Warn("Không thể kết nối đến Redis (%s:%s): %v. Cache sẽ tạm thời không khả dụng.", cfg.Host, cfg.Port, err)
	} else {
		utils.Info("Kết nối thành công đến Redis tại %s:%s!", cfg.Host, cfg.Port)
	}

	return rdb
}
