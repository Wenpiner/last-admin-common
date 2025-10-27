package userctx

import (
	"context"

	"google.golang.org/grpc/metadata"
)

// GetUserIDFromContext 从上下文中获取用户ID
func GetUserIDFromContext(ctx context.Context) (string, bool) {
	if userId, ok := ctx.Value("userId").(string); !ok {
		if md, ok := metadata.FromIncomingContext(ctx); !ok {
			return "", false
		} else {
			if data := md.Get("userId"); len(data) > 0 {
				return data[0], true
			}
		}
	} else {
		return userId, true
	}
	return "", false
}
