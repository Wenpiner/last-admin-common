package totp

// VerifyTotpCodeRequest TOTP验证请求
type VerifyTotpCodeRequest struct {
	UserId    string
	TotpCode  string
	IpAddress *string
	UserAgent *string
}

// VerifyTotpCodeResponse TOTP验证响应
type VerifyTotpCodeResponse struct {
	IsValid           bool
	Message           string
	RemainingAttempts *int32
	LockedUntil       *int64
}

