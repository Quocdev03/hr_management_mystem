package utils

import (
	"context"

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


