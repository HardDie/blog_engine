package utils

import (
	"context"

	"github.com/HardDie/blog_engine/internal/entity"
)

func GetUserIDFromContext(ctx context.Context) int64 {
	return ctx.Value("userID").(int64)
}
func GetSessionFromContext(ctx context.Context) *entity.Session {
	return ctx.Value("session").(*entity.Session)
}
