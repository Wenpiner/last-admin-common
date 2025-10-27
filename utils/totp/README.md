# TOTP 验证工具包

简洁的 TOTP（Time-based One-Time Password）验证工具包，用于验证用户的TOTP代码。

## 特性

- 简洁易用：只有一个验证函数
- 低耦合：不依赖业务逻辑，只负责验证
- 可复用：任何需要TOTP验证的模块都可以使用
- 易测试：可以Mock RPC客户端进行单元测试

## 安装

```bash
go get github.com/wenpiner/last-admin-common
```

## 快速开始

### 基础使用

```go
package main

import (
    "context"
    "github.com/wenpiner/last-admin-common/utils/totp"
)

func main() {
    // 创建验证器（需要注入UserService RPC客户端）
    verifier := totp.NewVerifier(userRpcClient)
    
    // 验证TOTP代码
    resp, err := verifier.Verify(ctx, userID, totpCode)
    if err != nil {
        // 处理错误
        return
    }
    
    if !resp.IsValid {
        // TOTP验证失败
        return
    }
    
    // TOTP验证成功，继续业务逻辑
}
```

### 带详细信息的验证

```go
// 验证TOTP代码并记录IP和User Agent
resp, err := verifier.VerifyWithDetails(
    ctx, 
    userID, 
    totpCode,
    ipAddress,
    userAgent,
)
```

## API 参考

### Verifier 方法

```go
// 创建验证器
func NewVerifier(userRpc UserServiceClient) *Verifier

// 验证TOTP代码
func (v *Verifier) Verify(ctx context.Context, userID, code string) (*VerifyTotpCodeResponse, error)

// 验证TOTP代码（带IP和UserAgent信息）
func (v *Verifier) VerifyWithDetails(ctx context.Context, userID, code, ipAddress, userAgent string) (*VerifyTotpCodeResponse, error)
```

### 数据结构

```go
// TOTP验证请求
type VerifyTotpCodeRequest struct {
    UserId    string
    TotpCode  string
    IpAddress *string
    UserAgent *string
}

// TOTP验证响应
type VerifyTotpCodeResponse struct {
    IsValid           bool
    Message           string
    RemainingAttempts *int32
    LockedUntil       *int64
}
```

## 使用场景

### 登录场景

```go
// 1. 验证用户名、密码、验证码
// 2. 检查用户是否启用TOTP
if user.TotpEnabled {
    // 返回用户ID作为TwoFactorKey
    return &LoginResponse{
        Data: LoginInfo{
            TwoFactorKey: user.Id,
        },
    }
}

// 3. 用户提交TOTP代码
resp, err := verifier.Verify(ctx, userID, totpCode)
if err != nil || !resp.IsValid {
    return nil, errors.New("TOTP验证失败")
}

// 4. 生成Token
accessToken := generateToken(userID, ...)
return &LoginResponse{
    Data: LoginInfo{
        AccessToken: accessToken,
    },
}
```

## 注意事项

1. **RPC客户端注入**：需要在初始化时注入UserService RPC客户端
2. **错误处理**：验证失败时会返回error，需要妥善处理
3. **响应检查**：即使没有error，也需要检查 `resp.IsValid` 字段
4. **账户锁定**：如果 `resp.LockedUntil` 不为nil，表示账户已被锁定

## 示例

完整的登录流程示例：

```go
// 第一步：登录
func (l *LoginLogic) Login(req *LoginRequest) (*LoginResponse, error) {
    // 验证验证码
    // 获取用户信息
    // 验证密码
    // 验证用户状态
    
    if user.TotpEnabled {
        return &LoginResponse{
            Data: LoginInfo{
                TwoFactorKey: user.Id,
            },
        }
    }
    
    // 生成Token
    token := generateToken(user.Id, ...)
    return &LoginResponse{
        Data: LoginInfo{
            AccessToken: token,
        },
    }
}

// 第二步：TOTP验证
func (l *TotpVerifyLogic) Verify(req *TotpVerifyRequest) (*TotpVerifyResponse, error) {
    // 验证TOTP代码
    resp, err := l.svcCtx.TotpVerifier.Verify(l.ctx, req.ID, req.Code)
    if err != nil {
        return nil, err
    }
    
    if !resp.IsValid {
        return nil, errors.New(resp.Message)
    }
    
    // 生成Token
    token := generateToken(req.ID, ...)
    return &TotpVerifyResponse{
        Data: LoginInfo{
            AccessToken: token,
        },
    }
}
```

