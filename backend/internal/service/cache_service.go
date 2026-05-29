package service

import (
	"chiquoc_hocgolang/internal/utils"
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/redis/go-redis/v9"
)

// CacheService cung cấp một lớp wrapper sạch sẽ để thao tác với Redis cache
type CacheService interface {
	Get(ctx context.Context, key string, dest interface{}) error
	Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error
	Delete(ctx context.Context, key string) error
	DeletePattern(ctx context.Context, pattern string) error
}

type cacheService struct {
	client *redis.Client
}

// NewCacheService tạo mới thực thể CacheService
func NewCacheService(client *redis.Client) CacheService {
	return &cacheService{
		client: client,
	}
}

func (c *cacheService) Get(ctx context.Context, key string, dest interface{}) error {
	if c.client == nil {
		return errors.New("redis client chưa được khởi tạo")
	}

	val, err := c.client.Get(ctx, key).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return errors.New("cache miss")
		}
		return err
	}

	err = json.Unmarshal([]byte(val), dest)
	if err != nil {
		utils.Error("Lỗi Unmarshal JSON cho cache key %s: %v", key, err)
		return err
	}

	return nil
}

func (c *cacheService) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	if c.client == nil {
		return errors.New("redis client chưa được khởi tạo")
	}

	bytes, err := json.Marshal(value)
	if err != nil {
		utils.Error("Lỗi Marshal JSON cho cache key %s: %v", key, err)
		return err
	}

	err = c.client.Set(ctx, key, bytes, expiration).Err()
	if err != nil {
		utils.Error("Lỗi Redis SET cho key %s: %v", key, err)
		return err
	}

	return nil
}

func (c *cacheService) Delete(ctx context.Context, key string) error {
	if c.client == nil {
		return errors.New("redis client chưa được khởi tạo")
	}

	err := c.client.Del(ctx, key).Err()
	if err != nil {
		utils.Error("Lỗi Redis DEL cho key %s: %v", key, err)
		return err
	}

	return nil
}

func (c *cacheService) DeletePattern(ctx context.Context, pattern string) error {
	if c.client == nil {
		return errors.New("redis client chưa được khởi tạo")
	}

	var cursor uint64
	for {
		keys, nextCursor, err := c.client.Scan(ctx, cursor, pattern, 100).Result()
		if err != nil {
			utils.Error("Lỗi Redis SCAN cho pattern %s: %v", pattern, err)
			return err
		}

		if len(keys) > 0 {
			err = c.client.Del(ctx, keys...).Err()
			if err != nil {
				utils.Error("Lỗi Redis DEL cho danh sách keys %v: %v", keys, err)
				return err
			}
		}

		cursor = nextCursor
		if cursor == 0 {
			break
		}
	}

	return nil
}
