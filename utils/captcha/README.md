# Captcha 验证码工具包

基于 `github.com/mojocn/base64Captcha` 的验证码生成和验证工具包，支持多种验证码类型和存储方式。

## 特性

- 支持多种验证码类型：数字、字符串、数学、中文、音频
- 支持多种存储方式：内存存储、Redis存储
- 提供全局服务和实例服务两种使用方式
- 完整的配置支持
- 线程安全

## 安装

```bash
go get github.com/mojocn/base64Captcha
go get github.com/redis/go-redis/v9
```

## 快速开始

### 使用默认配置

```go
package main

import (
    "fmt"
    "log"
    
    "github.com/wenpiner/last-admin-common/utils/captcha"
)

func main() {
    // 使用默认配置创建服务
    service, err := captcha.NewService(nil)
    if err != nil {
        log.Fatal(err)
    }
    defer service.Close()
    
    // 生成验证码
    result, err := service.Generate()
    if err != nil {
        log.Fatal(err)
    }
    
    fmt.Printf("验证码ID: %s\n", result.ID)
    fmt.Printf("验证码答案: %s\n", result.Answer)
    fmt.Printf("Base64图片: %s\n", result.Base64Blob)
    
    // 验证验证码
    isValid := service.VerifyAndClear(result.ID, result.Answer)
    fmt.Printf("验证结果: %t\n", isValid)
}
```

### 使用全局服务

```go
package main

import (
    "fmt"
    "log"
    
    "github.com/wenpiner/last-admin-common/utils/captcha"
)

func main() {
    // 初始化全局服务
    config := captcha.DefaultConfig()
    config.Store.Type = captcha.StoreTypeMemory
    
    err := captcha.InitGlobalService(config)
    if err != nil {
        log.Fatal(err)
    }
    defer captcha.CloseGlobalService()
    
    // 使用全局函数
    result, err := captcha.Generate()
    if err != nil {
        log.Fatal(err)
    }
    
    fmt.Printf("验证码ID: %s\n", result.ID)
    
    // 验证验证码
    isValid := captcha.VerifyAndClear(result.ID, result.Answer)
    fmt.Printf("验证结果: %t\n", isValid)
}
```

## 配置说明

### 验证码类型

```go
// 支持的验证码类型
const (
    TypeDigit   = "digit"   // 数字验证码
    TypeString  = "string"  // 字符串验证码
    TypeMath    = "math"    // 数学验证码
    TypeChinese = "chinese" // 中文验证码
    TypeAudio   = "audio"   // 音频验证码
    TypeRandom  = "random"  // 随机验证码类型
)
```

### 存储类型

```go
// 支持的存储类型
const (
    StoreTypeMemory = "memory" // 内存存储
    StoreTypeRedis  = "redis"  // Redis存储
)
```

### 自定义配置

```go
config := &captcha.CaptchaConfig{
    Type: captcha.TypeString,
    String: captcha.StringConfig{
        Height:     60,
        Width:      240,
        Length:     4,
        Source:     "abcdefghijklmnopqrstuvwxyz0123456789",
        NoiseCount: 0,
    },
    Store: captcha.StoreConfig{
        Type:      captcha.StoreTypeRedis,
        Expire:    5 * time.Minute,
        KeyPrefix: "captcha:",
        Redis: captcha.RedisConfig{
            Addr:     "localhost:6379",
            Password: "",
            DB:       0,
            PoolSize: 10,
        },
    },
}
```

### 随机验证码配置

```go
config := &captcha.CaptchaConfig{
    Type: captcha.TypeRandom,
    Random: captcha.RandomConfig{
        EnabledTypes: []captcha.CaptchaType{
            captcha.TypeDigit,
            captcha.TypeString,
            captcha.TypeMath,
        },
        ExcludeAudio: true, // 排除音频验证码
    },
    Store: captcha.StoreConfig{
        Type:      captcha.StoreTypeMemory,
        Expire:    5 * time.Minute,
        KeyPrefix: "captcha:",
    },
}
```

## API 参考

### Service 方法

```go
// 创建服务
func NewService(config *CaptchaConfig) (*Service, error)

// 生成不同类型的验证码
func (s *Service) Generate() (*CaptchaResult, error)
func (s *Service) GenerateDigit() (*CaptchaResult, error)
func (s *Service) GenerateString() (*CaptchaResult, error)
func (s *Service) GenerateMath() (*CaptchaResult, error)
func (s *Service) GenerateChinese() (*CaptchaResult, error)
func (s *Service) GenerateAudio() (*AudioCaptchaResult, error)
func (s *Service) GenerateRandom() (*CaptchaResult, error)

// 验证验证码
func (s *Service) Verify(id, answer string, clear bool) bool
func (s *Service) VerifyAndClear(id, answer string) bool

// 关闭服务
func (s *Service) Close() error
```

### 全局函数

```go
// 初始化和管理全局服务
func InitGlobalService(config *CaptchaConfig) error
func GetGlobalService() *Service
func CloseGlobalService() error

// 全局生成函数
func Generate() (*CaptchaResult, error)
func GenerateDigit() (*CaptchaResult, error)
func GenerateString() (*CaptchaResult, error)
func GenerateMath() (*CaptchaResult, error)
func GenerateChinese() (*CaptchaResult, error)
func GenerateAudio() (*AudioCaptchaResult, error)
func GenerateRandom() (*CaptchaResult, error)

// 全局验证函数
func Verify(id, answer string, clear bool) bool
func VerifyAndClear(id, answer string) bool
```

## 使用示例

### 生成不同类型的验证码

```go
service, _ := captcha.NewService(config)
defer service.Close()

// 数字验证码
digitResult, _ := service.GenerateDigit()

// 字符串验证码
stringResult, _ := service.GenerateString()

// 数学验证码
mathResult, _ := service.GenerateMath()

// 中文验证码
chineseResult, _ := service.GenerateChinese()

// 音频验证码
audioResult, _ := service.GenerateAudio()

// 随机验证码
randomResult, _ := service.GenerateRandom()
```

### 使用随机验证码

```go
// 配置随机验证码
config := captcha.DefaultConfig()
config.Type = captcha.TypeRandom
config.Random.EnabledTypes = []captcha.CaptchaType{
    captcha.TypeDigit,
    captcha.TypeString,
    captcha.TypeMath,
}

service, err := captcha.NewService(config)
if err != nil {
    log.Fatal(err)
}
defer service.Close()

// 生成随机类型的验证码
result, err := service.GenerateRandom()
if err != nil {
    log.Fatal(err)
}

fmt.Printf("验证码ID: %s\n", result.ID)
fmt.Printf("验证码答案: %s\n", result.Answer)
```

### 使用Redis存储

```go
config := captcha.DefaultConfig()
config.Store.Type = captcha.StoreTypeRedis
config.Store.Redis = captcha.RedisConfig{
    Addr:     "localhost:6379",
    Password: "your-password",
    DB:       0,
    PoolSize: 10,
}

service, err := captcha.NewService(config)
if err != nil {
    log.Fatal(err)
}
defer service.Close()
```

## 测试

运行测试：

```bash
go test -v ./utils/captcha
```

运行示例：

```bash
go test -v -run Example ./utils/captcha
```

## 注意事项

1. Redis存储需要确保Redis服务器可用
2. 音频验证码需要系统支持音频处理
3. 中文验证码需要相应的字体文件
4. 验证码有过期时间，默认5分钟
5. 建议在生产环境中使用Redis存储以支持分布式部署
