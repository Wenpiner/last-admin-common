package captcha

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/mojocn/base64Captcha"
)

// CaptchaResult 验证码生成结果
type CaptchaResult struct {
	ID          string `json:"id"`          // 验证码ID
	Base64Blob  string `json:"base64Blob"`  // Base64编码的图片数据
	Body        []byte `json:"body"`        // 原始图片数据
	Answer      string `json:"answer"`      // 验证码答案（仅用于测试）
	CaptchaType string `json:"captchaType"` // 验证码类型
}

// AudioCaptchaResult 音频验证码生成结果
type AudioCaptchaResult struct {
	ID         string `json:"id"`         // 验证码ID
	Base64Blob string `json:"base64Blob"` // Base64编码的音频数据
	Body       []byte `json:"body"`       // 原始音频数据
	Answer     string `json:"answer"`     // 验证码答案（仅用于测试）
}

// Service 验证码服务
type Service struct {
	config *CaptchaConfig
	store  Store
	driver *base64Captcha.Captcha
}

// NewService 创建验证码服务
func NewService(config *CaptchaConfig) (*Service, error) {
	if config == nil {
		config = DefaultConfig()
	}

	store, err := NewStore(config.Store)
	if err != nil {
		return nil, fmt.Errorf("failed to create store: %w", err)
	}

	driver := base64Captcha.NewCaptcha(nil, store)

	return &Service{
		config: config,
		store:  store,
		driver: driver,
	}, nil
}

// GenerateDigit 生成数字验证码
func (s *Service) GenerateDigit() (*CaptchaResult, error) {
	driver := base64Captcha.NewDriverDigit(
		s.config.Digit.Height,
		s.config.Digit.Width,
		s.config.Digit.Length,
		s.config.Digit.MaxSkew,
		s.config.Digit.DotCount,
	)

	captcha := base64Captcha.NewCaptcha(driver, s.store)
	id, b64s, answer, err := captcha.Generate()
	if err != nil {
		return nil, fmt.Errorf("failed to generate digit captcha: %w", err)
	}

	return &CaptchaResult{
		ID:         id,
		Base64Blob: b64s,
		Answer:     answer,
		CaptchaType: string(TypeDigit),
	}, nil
}

// GenerateString 生成字符串验证码
func (s *Service) GenerateString() (*CaptchaResult, error) {
	driver := base64Captcha.NewDriverString(
		s.config.String.Height,
		s.config.String.Width,
		s.config.String.NoiseCount,
		s.config.String.ShowLineOptions,
		s.config.String.Length,
		s.config.String.Source,
		s.config.String.BgColor,
		base64Captcha.DefaultEmbeddedFonts,
		s.config.String.Fonts,
	)

	captcha := base64Captcha.NewCaptcha(driver, s.store)
	id, b64s, answer, err := captcha.Generate()
	if err != nil {
		return nil, fmt.Errorf("failed to generate string captcha: %w", err)
	}

	return &CaptchaResult{
		ID:         id,
		Base64Blob: b64s,
		Answer:     answer,
		CaptchaType: string(TypeString),
	}, nil
}

// GenerateMath 生成数学验证码
func (s *Service) GenerateMath() (*CaptchaResult, error) {
	driver := base64Captcha.NewDriverMath(
		s.config.Math.Height,
		s.config.Math.Width,
		s.config.Math.NoiseCount,
		s.config.Math.ShowLineOptions,
		s.config.Math.BgColor,
		base64Captcha.DefaultEmbeddedFonts,
		s.config.Math.Fonts,
	)

	captcha := base64Captcha.NewCaptcha(driver, s.store)
	id, b64s, answer, err := captcha.Generate()
	if err != nil {
		return nil, fmt.Errorf("failed to generate math captcha: %w", err)
	}

	return &CaptchaResult{
		ID:         id,
		Base64Blob: b64s,
		Answer:     answer,
		CaptchaType: string(TypeMath),
	}, nil
}

// GenerateChinese 生成中文验证码
func (s *Service) GenerateChinese() (*CaptchaResult, error) {
	driver := base64Captcha.NewDriverChinese(
		s.config.Chinese.Height,
		s.config.Chinese.Width,
		s.config.Chinese.NoiseCount,
		s.config.Chinese.ShowLineOptions,
		s.config.Chinese.Length,
		s.config.Chinese.Source,
		s.config.Chinese.BgColor,
		base64Captcha.DefaultEmbeddedFonts,
		s.config.Chinese.Fonts,
	)

	captcha := base64Captcha.NewCaptcha(driver, s.store)
	id, b64s, answer, err := captcha.Generate()
	if err != nil {
		return nil, fmt.Errorf("failed to generate chinese captcha: %w", err)
	}

	return &CaptchaResult{
		ID:         id,
		Base64Blob: b64s,
		Answer:     answer,
		CaptchaType: string(TypeChinese),
	}, nil
}

// GenerateAudio 生成音频验证码
func (s *Service) GenerateAudio() (*AudioCaptchaResult, error) {
	driver := base64Captcha.NewDriverAudio(
		s.config.Audio.Length,
		s.config.Audio.Language,
	)

	captcha := base64Captcha.NewCaptcha(driver, s.store)
	id, b64s, answer, err := captcha.Generate()
	if err != nil {
		return nil, fmt.Errorf("failed to generate audio captcha: %w", err)
	}

	return &AudioCaptchaResult{
		ID:         id,
		Base64Blob: b64s,
		Answer:     answer,
	}, nil
}

// GenerateRandom 生成随机类型的验证码
func (s *Service) GenerateRandom() (*CaptchaResult, error) {
	randomType := s.getRandomCaptchaType()

	switch randomType {
	case TypeDigit:
		return s.GenerateDigit()
	case TypeString:
		return s.GenerateString()
	case TypeMath:
		return s.GenerateMath()
	case TypeChinese:
		return s.GenerateChinese()
	case TypeAudio:
		// 音频验证码返回特殊处理，这里转换为普通结果
		audioResult, err := s.GenerateAudio()
		if err != nil {
			return nil, err
		}
		return &CaptchaResult{
			ID:         audioResult.ID,
			Base64Blob: audioResult.Base64Blob,
			Answer:     audioResult.Answer,
		}, nil
	default:
		return s.GenerateDigit()
	}
}

// getRandomCaptchaType 根据配置随机选择一个验证码类型
func (s *Service) getRandomCaptchaType() CaptchaType {
	enabledTypes := s.config.Random.EnabledTypes

	// 如果没有配置启用类型，使用默认类型
	if len(enabledTypes) == 0 {
		enabledTypes = []CaptchaType{TypeDigit, TypeString, TypeMath, TypeChinese}
	}

	// 过滤掉音频类型（如果配置了排除音频）
	if s.config.Random.ExcludeAudio {
		var filteredTypes []CaptchaType
		for _, t := range enabledTypes {
			if t != TypeAudio {
				filteredTypes = append(filteredTypes, t)
			}
		}
		enabledTypes = filteredTypes
	}

	// 如果过滤后没有可用类型，使用数字验证码作为默认
	if len(enabledTypes) == 0 {
		return TypeDigit
	}

	rand.Seed(time.Now().UnixNano())
	return enabledTypes[rand.Intn(len(enabledTypes))]
}

// Generate 根据配置类型生成验证码
func (s *Service) Generate() (*CaptchaResult, error) {
	switch s.config.Type {
	case TypeDigit:
		return s.GenerateDigit()
	case TypeString:
		return s.GenerateString()
	case TypeMath:
		return s.GenerateMath()
	case TypeChinese:
		return s.GenerateChinese()
	case TypeRandom:
		return s.GenerateRandom()
	default:
		return s.GenerateDigit()
	}
}

// Verify 验证验证码
func (s *Service) Verify(id, answer string, clear bool) bool {
	return s.store.Verify(id, answer, clear)
}

// VerifyAndClear 验证验证码并清除
func (s *Service) VerifyAndClear(id, answer string) bool {
	return s.Verify(id, answer, true)
}

// Clear 清除验证码
func (s *Service) Clear(id string) {
	s.Clear(id)
}

// Close 关闭服务
func (s *Service) Close() error {
	if s.store != nil {
		return s.store.Close()
	}
	return nil
}
