package handler

import (
	"chiquoc_hocgolang/internal/utils"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type PerformanceHandler struct {
	db  *gorm.DB
	rdb *redis.Client
}

func NewPerformanceHandler(db *gorm.DB, rdb *redis.Client) *PerformanceHandler {
	return &PerformanceHandler{
		db:  db,
		rdb: rdb,
	}
}

// Generate large dummy data
func generateDummyData() []map[string]interface{} {
	var data []map[string]interface{}
	for i := 1; i <= 50; i++ {
		data = append(data, map[string]interface{}{
			"id":         i,
			"name":       fmt.Sprintf("Employee %d", i),
			"department": "Engineering",
			"salary":     1000 + i*10,
			"joined_at":  time.Now().AddDate(0, -i, 0).Format(time.RFC3339),
			"metadata":   fmt.Sprintf("Some random heavy metadata string to increase payload size. Iteration %d", i),
		})
	}
	return data
}

// TestPerformance handles GET /api/v1/performance/test
// query: ?use_cache=true|false
func (h *PerformanceHandler) TestPerformance(c *gin.Context) {
	useCache := c.Query("use_cache") == "true"
	cacheKey := "performance:test_data"
	redisCtx := c.Request.Context()

	start := time.Now()

	// 1. Nếu có dùng Cache, thử lấy từ Redis trước
	if useCache {
		val, err := h.rdb.Get(redisCtx, cacheKey).Result()
		if err == nil {
			// Cache Hit
			var result []map[string]interface{}
			_ = json.Unmarshal([]byte(val), &result)

			duration := time.Since(start).Milliseconds()

			c.JSON(http.StatusOK, gin.H{
				"success": true,
				"data": gin.H{
					"source":      "Redis Cache",
					"cache_hit":   true,
					"response_ms": duration,
					"items_count": len(result),
					"records":     result,
					"message":     "Dữ liệu được lấy ngay lập tức từ Redis (không chạm DB).",
				},
			})
			return
		} else if err != redis.Nil {
			utils.Warn("Lỗi Redis khi lấy cache: %v", err)
		}
	}

	// 2. Nếu không dùng Cache hoặc Cache Miss -> "Query DB"
	// Giả lập một câu query phức tạp mất khoảng 500ms
	time.Sleep(500 * time.Millisecond)

	// Có thể query thực tế:
	// h.db.Raw("SELECT SLEEP(0.5)").Scan(&var) // (MySQL specific)

	data := generateDummyData()

	// 3. Nếu đang dùng chế độ Cache nhưng bị Miss, lưu lại vào Redis
	if useCache {
		bytesData, _ := json.Marshal(data)
		err := h.rdb.Set(redisCtx, cacheKey, bytesData, 1*time.Hour).Err()
		if err != nil {
			utils.Warn("Không thể lưu Redis cache: %v", err)
		}
	}

	duration := time.Since(start).Milliseconds()

	sourceStr := "MySQL Database"
	cacheHit := false
	msg := "Dữ liệu được lấy trực tiếp từ Database (Không dùng Cache)."
	if useCache {
		msg = "Cache Miss: Đã query Database và lưu kết quả vào Redis cho lần sau."
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data": gin.H{
			"source":      sourceStr,
			"cache_hit":   cacheHit,
			"response_ms": duration,
			"items_count": len(data),
			"records":     data,
			"message":     msg,
		},
	})
}

// ClearCache handles DELETE /api/v1/performance/clear
func (h *PerformanceHandler) ClearCache(c *gin.Context) {
	cacheKey := "performance:test_data"
	err := h.rdb.Del(c.Request.Context(), cacheKey).Err()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   gin.H{"code": 500, "message": "Lỗi khi xóa cache"},
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Đã xóa cache thành công. API tiếp theo sẽ truy vấn lại vào Database (Cold Start).",
	})
}
