package shared

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"
)

func DashboardCacheKey(userID string) string {
	return fmt.Sprintf("dashboard:%s", userID)
}

func InvalidateDashboard(ctx context.Context, redisClient *redis.Client, userID string) {
	if redisClient == nil {
		return
	}
	_ = redisClient.Del(ctx, DashboardCacheKey(userID)).Err()
}
