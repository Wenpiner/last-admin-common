package captcha

import (
	"fmt"
	"log"
	"testing"
	"time"
)

// ExampleService_GenerateDigit 演示如何生成数字验证码
func ExampleService_GenerateDigit() {
	// 创建配置
	config := DefaultConfig()
	config.Store.Type = StoreTypeMemory // 使用内存存储用于示例

	// 创建服务
	service, err := NewService(config)
	if err != nil {
		log.Fatal(err)
	}
	defer service.Close()

	// 生成数字验证码
	result, err := service.GenerateDigit()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("验证码ID: %s\n", result.ID)
	fmt.Printf("验证码答案: %s\n", result.Answer)
	fmt.Printf("Base64图片数据长度: %d\n", len(result.Base64Blob))

	// 验证验证码
	isValid := service.VerifyAndClear(result.ID, result.Answer)
	fmt.Printf("验证结果: %t\n", isValid)
}

// ExampleInitGlobalService 演示如何使用全局服务
func ExampleInitGlobalService() {
	// 创建配置
	config := &CaptchaConfig{
		Type: TypeString,
		String: StringConfig{
			Height:     60,
			Width:      240,
			Length:     4,
			Source:     "abcdefghijklmnopqrstuvwxyz0123456789",
			NoiseCount: 0,
		},
		Store: StoreConfig{
			Type:      StoreTypeMemory,
			Expire:    5 * time.Minute,
			KeyPrefix: "captcha:",
		},
	}

	// 初始化全局服务
	err := InitGlobalService(config)
	if err != nil {
		log.Fatal(err)
	}
	defer CloseGlobalService()

	// 使用全局函数生成验证码
	result, err := Generate()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("验证码ID: %s\n", result.ID)
	fmt.Printf("验证码答案: %s\n", result.Answer)

	// 验证验证码
	isValid := VerifyAndClear(result.ID, result.Answer)
	fmt.Printf("验证结果: %t\n", isValid)
}

// ExampleRedisStore 演示如何使用Redis存储
func ExampleRedisStore() {
	// 注意：这个示例需要Redis服务器运行在localhost:6379
	config := &CaptchaConfig{
		Type: TypeMath,
		Store: StoreConfig{
			Type:      StoreTypeRedis,
			Expire:    10 * time.Minute,
			KeyPrefix: "captcha:",
			Redis: RedisConfig{
				Addr:     "localhost:6379",
				Password: "",
				DB:       0,
				PoolSize: 10,
			},
		},
	}

	service, err := NewService(config)
	if err != nil {
		log.Printf("Redis连接失败，跳过示例: %v", err)
		return
	}
	defer service.Close()

	// 生成数学验证码
	result, err := service.GenerateMath()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("验证码ID: %s\n", result.ID)
	fmt.Printf("验证码答案: %s\n", result.Answer)

	// 验证验证码
	isValid := service.VerifyAndClear(result.ID, result.Answer)
	fmt.Printf("验证结果: %t\n", isValid)
}

// TestExampleMultipleTypes 演示不同类型的验证码
func TestExampleMultipleTypes(t *testing.T) {
	config := DefaultConfig()
	config.Store.Type = StoreTypeMemory

	service, err := NewService(config)
	if err != nil {
		log.Fatal(err)
	}
	defer service.Close()

	// 生成不同类型的验证码
	types := []struct {
		name string
		fn   func() (*CaptchaResult, error)
	}{
		{"数字验证码", service.GenerateDigit},
		{"字符串验证码", service.GenerateString},
		{"数学验证码", service.GenerateMath},
		{"中文验证码", service.GenerateChinese},
	}

	for _, typ := range types {
		result, err := typ.fn()
		if err != nil {
			t.Logf("%s生成失败: %v", typ.name, err)
			continue
		}

		t.Logf("%s - ID: %s, 答案: %s\n", typ.name, result.ID, result.Answer)

		// 验证
		isValid := service.Verify(result.ID, result.Answer, true)
		t.Logf("%s验证结果: %t\n", typ.name, isValid)
	}
}

// TestExampleAudioCaptcha 演示音频验证码
func TestExampleAudioCaptcha(t *testing.T) {
	config := DefaultConfig()
	config.Store.Type = StoreTypeMemory

	service, err := NewService(config)
	if err != nil {
		log.Fatal(err)
	}
	defer service.Close()

	// 生成音频验证码
	result, err := service.GenerateAudio()
	if err != nil {
		log.Fatal(err)
	}

	t.Logf("音频验证码ID: %s\n", result.ID)
	t.Logf("音频验证码答案: %s\n", result.Answer)
	t.Logf("Base64音频数据长度: %d\n", len(result.Base64Blob))

	// 验证验证码
	isValid := service.VerifyAndClear(result.ID, result.Answer)
	t.Logf("验证结果: %t\n", isValid)
}

// TestExampleRandomCaptcha 演示随机验证码
func TestExampleRandomCaptcha(t *testing.T) {
	config := DefaultConfig()
	config.Store.Type = StoreTypeMemory
	config.Type = TypeRandom

	// 自定义随机配置
	config.Random.EnabledTypes = []CaptchaType{TypeDigit, TypeString, TypeMath}
	config.Random.ExcludeAudio = true

	service, err := NewService(config)
	if err != nil {
		log.Fatal(err)
	}
	defer service.Close()

	// 生成多个随机验证码
	for i := 0; i < 5; i++ {
		result, err := service.GenerateRandom()
		if err != nil {
			log.Fatal(err)
		}

		t.Logf("随机验证码 %d - ID: %s, 答案: %s\n", i+1, result.ID, result.Answer)

		// 验证验证码
		isValid := service.VerifyAndClear(result.ID, result.Answer)
		t.Logf("验证结果: %t\n", isValid)
	}
}

// TestExampleGenerateRandom 演示全局随机验证码服务
func TestExampleGenerateRandom(t *testing.T) {
	// 创建随机验证码配置
	config := DefaultConfig()
	config.Type = TypeRandom
	config.Store.Type = StoreTypeMemory
	config.Random.EnabledTypes = []CaptchaType{TypeDigit, TypeString, TypeMath, TypeChinese}
	config.Random.ExcludeAudio = true

	// 初始化全局服务
	err := InitGlobalService(config)
	if err != nil {
		log.Fatal(err)
	}
	defer CloseGlobalService()

	// 使用全局函数生成随机验证码
	result, err := GenerateRandom()
	if err != nil {
		log.Fatal(err)
	}

	t.Logf("全局随机验证码ID: %s\n", result.ID)
	t.Logf("全局随机验证码答案: %s\n", result.Answer)

	// 验证验证码
	isValid := VerifyAndClear(result.ID, result.Answer)
	t.Logf("验证结果: %t\n", isValid)
}
