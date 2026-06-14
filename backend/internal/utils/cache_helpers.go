package utils

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"
)

// InvalidateDashboardStats xóa cache thống kê dashboard.
// Gọi sau khi tạo/cập nhật/xóa employee, department, hoặc user.
// Graceful fail-open: nếu rdb == nil thì bỏ qua (không crash).
func InvalidateDashboardStats(ctx context.Context, rdb *redis.Client) error {
	if rdb == nil {
		return nil
	}

	if err := rdb.Del(ctx, "dashboard:stats").Err(); err != nil {
		Error("Không thể xóa dashboard stats cache: %v", err)
		return err
	}

	Info("Dashboard stats cache đã được xóa")
	return nil
}

// InvalidateRefreshToken xóa refresh token của một user khỏi Redis.
// Graceful fail-open: nếu rdb == nil thì bỏ qua.
func InvalidateRefreshToken(ctx context.Context, rdb *redis.Client, userID uint) error {
	if rdb == nil {
		return nil
	}

	key := fmt.Sprintf("refresh_token:%d", userID)
	return rdb.Del(ctx, key).Err()
}
