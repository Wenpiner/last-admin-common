package rolectx

import (
	"context"
	"strings"

	"github.com/zeromicro/go-zero/rest/enums"
	"google.golang.org/grpc/metadata"
)

func GetRoleFromContext(ctx context.Context) ([]string, bool) {
	if role, ok := ctx.Value("roleId").(string); !ok {
		if md, ok := metadata.FromIncomingContext(ctx); !ok {
			return nil, false
		} else {
			if data := md.Get(string(enums.RoleKey)); len(data) > 0 {
				roleIds := strings.Split(data[0], ",")
				return roleIds, true
			}
			return nil, false
		}
	}else{
		roleIds := strings.Split(role, ",")
		return roleIds, true
	}
}
