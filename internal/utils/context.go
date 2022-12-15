package utils

import "context"

func GetUserIDFromContext(ctx context.Context) int32 {
	return ctx.Value("userID").(int32)
}
