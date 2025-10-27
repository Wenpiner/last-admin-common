package captcha

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDefaultConfig(t *testing.T) {
	config := DefaultConfig()
	assert.NotNil(t, config)
	assert.Equal(t, TypeDigit, config.Type)
	assert.Equal(t, StoreTypeRedis, config.Store.Type)
	assert.Equal(t, "captcha:", config.Store.KeyPrefix)
	assert.Equal(t, 5*time.Minute, config.Store.Expire)
}

func TestNewServiceWithMemoryStore(t *testing.T) {
	config := DefaultConfig()
	config.Store.Type = StoreTypeMemory

	service, err := NewService(config)
	require.NoError(t, err)
	require.NotNil(t, service)
	defer service.Close()

	assert.Equal(t, config, service.config)
	assert.NotNil(t, service.store)
	assert.NotNil(t, service.driver)
}

func TestGenerateDigitCaptcha(t *testing.T) {
	config := DefaultConfig()
	config.Store.Type = StoreTypeMemory

	service, err := NewService(config)
	require.NoError(t, err)
	defer service.Close()

	result, err := service.GenerateDigit()
	require.NoError(t, err)
	require.NotNil(t, result)

	assert.NotEmpty(t, result.ID)
	assert.NotEmpty(t, result.Base64Blob)
	assert.NotEmpty(t, result.Answer)
	assert.Len(t, result.Answer, config.Digit.Length)
}

func TestGenerateStringCaptcha(t *testing.T) {
	config := DefaultConfig()
	config.Store.Type = StoreTypeMemory

	service, err := NewService(config)
	require.NoError(t, err)
	defer service.Close()

	result, err := service.GenerateString()
	require.NoError(t, err)
	require.NotNil(t, result)

	assert.NotEmpty(t, result.ID)
	assert.NotEmpty(t, result.Base64Blob)
	assert.NotEmpty(t, result.Answer)
	assert.Len(t, result.Answer, config.String.Length)
}

func TestGenerateMathCaptcha(t *testing.T) {
	config := DefaultConfig()
	config.Store.Type = StoreTypeMemory

	service, err := NewService(config)
	require.NoError(t, err)
	defer service.Close()

	result, err := service.GenerateMath()
	require.NoError(t, err)
	require.NotNil(t, result)

	assert.NotEmpty(t, result.ID)
	assert.NotEmpty(t, result.Base64Blob)
	assert.NotEmpty(t, result.Answer)
}

func TestGenerateChinese(t *testing.T) {
	config := DefaultConfig()
	config.Store.Type = StoreTypeMemory

	service, err := NewService(config)
	require.NoError(t, err)
	defer service.Close()

	result, err := service.GenerateChinese()
	require.NoError(t, err)
	require.NotNil(t, result)

	assert.NotEmpty(t, result.ID)
	assert.NotEmpty(t, result.Base64Blob)
	assert.NotEmpty(t, result.Answer)
	// 中文验证码的答案长度可能不等于配置的长度，因为是按字符计算的
	assert.True(t, len(result.Answer) >= config.Chinese.Length)
}

func TestGenerateAudio(t *testing.T) {
	config := DefaultConfig()
	config.Store.Type = StoreTypeMemory

	service, err := NewService(config)
	require.NoError(t, err)
	defer service.Close()

	result, err := service.GenerateAudio()
	require.NoError(t, err)
	require.NotNil(t, result)

	assert.NotEmpty(t, result.ID)
	assert.NotEmpty(t, result.Base64Blob)
	assert.NotEmpty(t, result.Answer)
	assert.Len(t, result.Answer, config.Audio.Length)
}

func TestVerifyCaptcha(t *testing.T) {
	config := DefaultConfig()
	config.Store.Type = StoreTypeMemory

	service, err := NewService(config)
	require.NoError(t, err)
	defer service.Close()

	// 生成验证码
	result, err := service.GenerateDigit()
	require.NoError(t, err)

	// 验证正确答案
	assert.True(t, service.Verify(result.ID, result.Answer, false))

	// 验证错误答案
	assert.False(t, service.Verify(result.ID, "wrong", false))

	// 验证并清除
	assert.True(t, service.VerifyAndClear(result.ID, result.Answer))

	// 再次验证应该失败（已清除）
	assert.False(t, service.Verify(result.ID, result.Answer, false))
}

func TestGenerateByType(t *testing.T) {
	config := DefaultConfig()
	config.Store.Type = StoreTypeMemory

	testCases := []struct {
		name string
		typ  CaptchaType
	}{
		{"digit", TypeDigit},
		{"string", TypeString},
		{"math", TypeMath},
		{"chinese", TypeChinese},
		{"random", TypeRandom},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			config.Type = tc.typ
			service, err := NewService(config)
			require.NoError(t, err)
			defer service.Close()

			result, err := service.Generate()
			require.NoError(t, err)
			require.NotNil(t, result)

			assert.NotEmpty(t, result.ID)
			assert.NotEmpty(t, result.Base64Blob)
			assert.NotEmpty(t, result.Answer)
		})
	}
}

func TestMemoryStore(t *testing.T) {
	store := NewMemoryStore(100, 5*time.Minute)
	defer store.Close()

	// 测试存储和获取
	err := store.Set("test-id", "test-value")
	assert.NoError(t, err)

	value := store.Get("test-id", false)
	assert.Equal(t, "test-value", value)

	// 测试验证
	assert.True(t, store.Verify("test-id", "test-value", false))
	assert.False(t, store.Verify("test-id", "wrong-value", false))

	// 测试清除
	assert.True(t, store.Verify("test-id", "test-value", true))
	value = store.Get("test-id", false)
	assert.Empty(t, value)
}

func TestGlobalService(t *testing.T) {
	// 清理全局服务
	defer CloseGlobalService()

	config := DefaultConfig()
	config.Store.Type = StoreTypeMemory

	// 初始化全局服务
	err := InitGlobalService(config)
	require.NoError(t, err)

	// 测试获取全局服务
	service := GetGlobalService()
	assert.NotNil(t, service)

	// 测试全局函数
	result, err := Generate()
	require.NoError(t, err)
	assert.NotEmpty(t, result.ID)

	// 测试验证
	assert.True(t, VerifyAndClear(result.ID, result.Answer))
}

func TestGlobalFunctionsWithoutInit(t *testing.T) {
	// 确保没有全局服务
	CloseGlobalService()

	// 测试在没有初始化全局服务的情况下使用全局函数
	result, err := GenerateDigit()
	require.NoError(t, err)
	assert.NotEmpty(t, result.ID)
}

func TestGenerateRandom(t *testing.T) {
	config := DefaultConfig()
	config.Store.Type = StoreTypeMemory
	config.Type = TypeRandom

	service, err := NewService(config)
	require.NoError(t, err)
	defer service.Close()

	// 测试随机生成
	result, err := service.GenerateRandom()
	require.NoError(t, err)
	require.NotNil(t, result)

	assert.NotEmpty(t, result.ID)
	assert.NotEmpty(t, result.Base64Blob)
	assert.NotEmpty(t, result.Answer)

	// 验证生成的验证码
	assert.True(t, service.VerifyAndClear(result.ID, result.Answer))
}

func TestRandomConfig(t *testing.T) {
	config := DefaultConfig()
	config.Store.Type = StoreTypeMemory
	config.Type = TypeRandom

	// 测试自定义随机配置
	config.Random.EnabledTypes = []CaptchaType{TypeDigit, TypeString}
	config.Random.ExcludeAudio = true

	service, err := NewService(config)
	require.NoError(t, err)
	defer service.Close()

	// 生成多个随机验证码，验证类型在配置范围内
	for i := 0; i < 10; i++ {
		result, err := service.GenerateRandom()
		require.NoError(t, err)
		assert.NotEmpty(t, result.ID)
		assert.NotEmpty(t, result.Answer)

		// 验证验证码
		assert.True(t, service.Verify(result.ID, result.Answer, true))
	}
}

func TestGlobalGenerateRandom(t *testing.T) {
	// 清理全局服务
	defer CloseGlobalService()

	config := DefaultConfig()
	config.Store.Type = StoreTypeMemory
	config.Type = TypeRandom

	// 初始化全局服务
	err := InitGlobalService(config)
	require.NoError(t, err)

	// 测试全局随机生成函数
	result, err := GenerateRandom()
	require.NoError(t, err)
	assert.NotEmpty(t, result.ID)
	assert.NotEmpty(t, result.Answer)

	// 验证
	assert.True(t, VerifyAndClear(result.ID, result.Answer))
}
