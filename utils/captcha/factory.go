package captcha

import (
	"sync"
)

var (
	// 全局服务实例
	globalService *Service
	globalMutex   sync.RWMutex
)

// InitGlobalService 初始化全局验证码服务
func InitGlobalService(config *CaptchaConfig) error {
	globalMutex.Lock()
	defer globalMutex.Unlock()

	if globalService != nil {
		globalService.Close()
	}

	service, err := NewService(config)
	if err != nil {
		return err
	}

	globalService = service
	return nil
}

// GetGlobalService 获取全局验证码服务
func GetGlobalService() *Service {
	globalMutex.RLock()
	defer globalMutex.RUnlock()
	return globalService
}

// Generate 使用全局服务生成验证码
func Generate() (*CaptchaResult, error) {
	service := GetGlobalService()
	if service == nil {
		// 如果没有初始化全局服务，使用默认配置创建临时服务
		service, err := NewService(nil)
		if err != nil {
			return nil, err
		}
		defer service.Close()
		return service.Generate()
	}
	return service.Generate()
}

// GenerateDigit 使用全局服务生成数字验证码
func GenerateDigit() (*CaptchaResult, error) {
	service := GetGlobalService()
	if service == nil {
		service, err := NewService(nil)
		if err != nil {
			return nil, err
		}
		defer service.Close()
		return service.GenerateDigit()
	}
	return service.GenerateDigit()
}

// GenerateString 使用全局服务生成字符串验证码
func GenerateString() (*CaptchaResult, error) {
	service := GetGlobalService()
	if service == nil {
		service, err := NewService(nil)
		if err != nil {
			return nil, err
		}
		defer service.Close()
		return service.GenerateString()
	}
	return service.GenerateString()
}

// GenerateMath 使用全局服务生成数学验证码
func GenerateMath() (*CaptchaResult, error) {
	service := GetGlobalService()
	if service == nil {
		service, err := NewService(nil)
		if err != nil {
			return nil, err
		}
		defer service.Close()
		return service.GenerateMath()
	}
	return service.GenerateMath()
}

// GenerateChinese 使用全局服务生成中文验证码
func GenerateChinese() (*CaptchaResult, error) {
	service := GetGlobalService()
	if service == nil {
		service, err := NewService(nil)
		if err != nil {
			return nil, err
		}
		defer service.Close()
		return service.GenerateChinese()
	}
	return service.GenerateChinese()
}

// GenerateAudio 使用全局服务生成音频验证码
func GenerateAudio() (*AudioCaptchaResult, error) {
	service := GetGlobalService()
	if service == nil {
		service, err := NewService(nil)
		if err != nil {
			return nil, err
		}
		defer service.Close()
		return service.GenerateAudio()
	}
	return service.GenerateAudio()
}

// GenerateRandom 使用全局服务生成随机类型验证码
func GenerateRandom() (*CaptchaResult, error) {
	service := GetGlobalService()
	if service == nil {
		service, err := NewService(nil)
		if err != nil {
			return nil, err
		}
		defer service.Close()
		return service.GenerateRandom()
	}
	return service.GenerateRandom()
}

// Verify 使用全局服务验证验证码
func Verify(id, answer string, clear bool) bool {
	service := GetGlobalService()
	if service == nil {
		return false
	}
	return service.Verify(id, answer, clear)
}

// VerifyAndClear 使用全局服务验证验证码并清除
func VerifyAndClear(id, answer string) bool {
	return Verify(id, answer, true)
}

// CloseGlobalService 关闭全局服务
func CloseGlobalService() error {
	globalMutex.Lock()
	defer globalMutex.Unlock()

	if globalService != nil {
		err := globalService.Close()
		globalService = nil
		return err
	}
	return nil
}
