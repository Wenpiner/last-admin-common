package totp

import (
	"context"
)

// UserServiceClient TOTP验证所需的RPC客户端接口
type UserServiceClient interface {
	// VerifyTotpCode 验证TOTP代码
	VerifyTotpCode(ctx context.Context, req *VerifyTotpCodeRequest) (*VerifyTotpCodeResponse, error)
}

// Verifier TOTP验证器
type Verifier struct {
	userRpc UserServiceClient
}

// NewVerifier 创建TOTP验证器
// userRpc: UserService RPC客户端
func NewVerifier(userRpc UserServiceClient) *Verifier {
	return &Verifier{
		userRpc: userRpc,
	}
}

// Verify 验证TOTP代码
// ctx: 上下文
// userID: 用户ID
// code: TOTP代码
// 返回: 验证响应，如果验证失败返回error
func (v *Verifier) Verify(ctx context.Context, userID, code string) (*VerifyTotpCodeResponse, error) {
	req := &VerifyTotpCodeRequest{
		UserId:   userID,
		TotpCode: code,
	}
	return v.userRpc.VerifyTotpCode(ctx, req)
}

// VerifyWithDetails 验证TOTP代码（带IP和UserAgent信息）
// ctx: 上下文
// userID: 用户ID
// code: TOTP代码
// ipAddress: IP地址
// userAgent: User Agent
// 返回: 验证响应，如果验证失败返回error
func (v *Verifier) VerifyWithDetails(ctx context.Context, userID, code, ipAddress, userAgent string) (*VerifyTotpCodeResponse, error) {
	req := &VerifyTotpCodeRequest{
		UserId:    userID,
		TotpCode:  code,
		IpAddress: &ipAddress,
		UserAgent: &userAgent,
	}
	return v.userRpc.VerifyTotpCode(ctx, req)
}

