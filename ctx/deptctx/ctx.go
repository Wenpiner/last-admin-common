package deptctx

import (
	"context"
	"strconv"

	"github.com/zeromicro/go-zero/rest/enums"
	"google.golang.org/grpc/metadata"
)

func GetDeptIDFromContext(ctx context.Context) (uint32, bool) {
	if deptId, ok := ctx.Value(enums.DeptKey).(uint32); !ok {
		if md, ok := metadata.FromIncomingContext(ctx); !ok {
			return 0, false
		} else {
			if data := md.Get(string(enums.DeptKey)); len(data) > 0 {
				deptId, err := strconv.Atoi(data[0])
				if err != nil {
					return 0, false
				}
				return uint32(deptId), true
			}
			return 0, false
		}
	} else {
		return deptId, true
	}
}
