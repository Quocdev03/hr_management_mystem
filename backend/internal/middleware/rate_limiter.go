package middleware

import (
	"chiquoc_hocgolang/internal/utils"
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

// RateLimiter áp dụng giới hạn tần suất yêu cầu sử dụng Redis (Fixed Window Algorithm)
func RateLimiter(rdb *redis.Client, limit int, window time.Duration) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		if rdb == nil {
			// Fail-open: Nếu Redis chưa được cấu hình, bỏ qua rate limit
			ctx.Next()
			return
		}

		ip := ctx.ClientIP()
		key := fmt.Sprintf("rate_limit:login:%s", ip)
		redisCtx := ctx.Request.Context()

		// Thực hiện tăng bộ đếm lên 1
		count, err := rdb.Incr(redisCtx, key).Result()
		if err != nil {
			// Fail-open: Nếu lỗi kết nối Redis, ghi log và cho qua
			utils.Error("Lỗi kết nối Redis trong middleware RateLimiter cho IP %s: %v", ip, err)
			ctx.Next()
			return
		}

		// Nếu là lượt truy cập đầu tiên trong cửa sổ thời gian, đặt thời gian hết hạn (TTL)
		if count == 1 {
			err = rdb.Expire(redisCtx, key, window).Err()
			if err != nil {
				utils.Error("Lỗi khi thiết lập TTL cho rate limit key %s: %v", key, err)
			}
		}

		// Tính toán số giây còn lại cho cửa sổ thời gian
		ttl, _ := rdb.TTL(redisCtx, key).Result()
		retryAfterSeconds := int(ttl.Seconds())
		if retryAfterSeconds < 0 {
			retryAfterSeconds = int(window.Seconds())
		}

		// Thêm các thông tin rate limit vào Response Header (chuẩn REST API)
		ctx.Header("X-RateLimit-Limit", fmt.Sprintf("%d", limit))
		ctx.Header("X-RateLimit-Remaining", fmt.Sprintf("%d", maxVal(0, int64(limit)-count)))
		ctx.Header("X-RateLimit-Reset", fmt.Sprintf("%d", retryAfterSeconds))

		// Nếu vượt quá giới hạn
		if count > int64(limit) {
			utils.Warn("IP %s đã vượt quá giới hạn lượt đăng nhập (%d/%d)", ip, count, limit)
			utils.TooManyRequests(ctx, fmt.Sprintf("Bạn đã thử đăng nhập quá nhiều lần! Vui lòng thử lại sau %d giây.", retryAfterSeconds))
			ctx.Abort()
			return
		}

		ctx.Next()
	}
}

func maxVal(a, b int64) int64 {
	if a > b {
		return a
	}
	return b
}
