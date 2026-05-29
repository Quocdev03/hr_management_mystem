package middleware

import (
	"bytes"
	"chiquoc_hocgolang/internal/utils"
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

// bodyLogWriter được dùng để chặn (intercept) và lưu lại kết quả trả về của API
type bodyLogWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w bodyLogWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

// CacheResponse là middleware tự động lưu cache các API GET vào Redis
func CacheResponse(rdb *redis.Client, expiration time.Duration) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// Chỉ cache các request GET
		if ctx.Request.Method != http.MethodGet {
			ctx.Next()
			return
		}

		// Tạo key cache dựa trên URL (bao gồm cả query params, vd: ?page=1)
		key := "cache:" + ctx.Request.URL.RequestURI()
		redisCtx := context.Background()

		// 1. Kiểm tra xem Redis có lưu cache cho key này chưa
		val, err := rdb.Get(redisCtx, key).Result()
		if err == nil {
			// Cache Hit: Trả về kết quả JSON trực tiếp từ Redis
			utils.Info("Cache hit cho URL: %s", ctx.Request.URL.RequestURI())
			ctx.Data(http.StatusOK, "application/json; charset=utf-8", []byte(val))
			ctx.Abort() // Ngăn không cho chạy xuống Handler/DB
			return
		}

		// 2. Cache Miss: Cho phép request đi tiếp và chặn kết quả trả về
		blw := &bodyLogWriter{body: bytes.NewBufferString(""), ResponseWriter: ctx.Writer}
		ctx.Writer = blw

		ctx.Next()

		// 3. Sau khi Handler xử lý xong, nếu status là 200 OK thì lưu kết quả vào Redis
		if ctx.Writer.Status() == http.StatusOK {
			err := rdb.Set(redisCtx, key, blw.body.String(), expiration).Err()
			if err != nil {
				utils.Error("Không thể lưu cache cho key %s: %v", key, err)
			} else {
				utils.Info("Đã lưu cache cho URL: %s (TTL: %v)", ctx.Request.URL.RequestURI(), expiration)
			}
		}
	}
}

// ClearCache là middleware dùng cho POST/PUT/DELETE để xoá cache cũ
func ClearCache(rdb *redis.Client, pattern string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// Cho phép request POST/PUT/DELETE chạy xuống DB trước
		ctx.Next()

		// Nếu xử lý thành công (status 2xx) thì tiến hành xoá cache
		status := ctx.Writer.Status()
		if status >= http.StatusOK && status < http.StatusMultipleChoices {
			redisCtx := context.Background()
			
			// Dùng SCAN để tìm tất cả các key match với pattern (thay vì KEYS để tránh block server)
			var cursor uint64
			for {
				var keys []string
				var err error
				keys, cursor, err = rdb.Scan(redisCtx, cursor, pattern, 100).Result()
				if err != nil {
					utils.Error("Lỗi khi tìm key cache pattern %s: %v", pattern, err)
					break
				}

				if len(keys) > 0 {
					err = rdb.Del(redisCtx, keys...).Err()
					if err != nil {
						utils.Error("Lỗi khi xoá các key cache %v: %v", keys, err)
					} else {
						utils.Info("Đã xoá %d key cache cũ cho pattern: %s", len(keys), pattern)
					}
				}

				if cursor == 0 {
					break
				}
			}
		}
	}
}
