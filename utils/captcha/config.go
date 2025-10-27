package captcha

import (
	"image/color"
	"time"
)

// CaptchaConfig 验证码配置
type CaptchaConfig struct {
	// 验证码类型配置
	Type CaptchaType `json:"type" yaml:"type"`
	
	// 数字验证码配置
	Digit DigitConfig `json:"digit" yaml:"digit"`
	
	// 字符串验证码配置
	String StringConfig `json:"string" yaml:"string"`
	
	// 数学验证码配置
	Math MathConfig `json:"math" yaml:"math"`
	
	// 中文验证码配置
	Chinese ChineseConfig `json:"chinese" yaml:"chinese"`
	
	// 音频验证码配置
	Audio AudioConfig `json:"audio" yaml:"audio"`
	
	// 存储配置
	Store StoreConfig `json:"store" yaml:"store"`

	// 随机验证码配置
	Random RandomConfig `json:"random" yaml:"random"`
}

// CaptchaType 验证码类型
type CaptchaType string

const (
	TypeDigit   CaptchaType = "digit"   // 数字验证码
	TypeString  CaptchaType = "string"  // 字符串验证码
	TypeMath    CaptchaType = "math"    // 数学验证码
	TypeChinese CaptchaType = "chinese" // 中文验证码
	TypeAudio   CaptchaType = "audio"   // 音频验证码
	TypeRandom  CaptchaType = "random"  // 随机验证码类型
)

// DigitConfig 数字验证码配置
type DigitConfig struct {
	Height   int `json:"height" yaml:"height"`     // 图片高度
	Width    int `json:"width" yaml:"width"`       // 图片宽度
	Length   int `json:"length" yaml:"length"`     // 验证码长度
	MaxSkew  float64 `json:"maxSkew" yaml:"maxSkew"`   // 最大倾斜角度
	DotCount int `json:"dotCount" yaml:"dotCount"` // 干扰点数量
}

// StringConfig 字符串验证码配置
type StringConfig struct {
	Height      int    `json:"height" yaml:"height"`           // 图片高度
	Width       int    `json:"width" yaml:"width"`             // 图片宽度
	NoiseCount  int    `json:"noiseCount" yaml:"noiseCount"`   // 干扰线数量
	ShowLineOptions int `json:"showLineOptions" yaml:"showLineOptions"` // 显示线条选项
	Length      int    `json:"length" yaml:"length"`           // 验证码长度
	Source      string `json:"source" yaml:"source"`           // 字符源
	BgColor     *color.RGBA `json:"bgColor" yaml:"bgColor"` // 背景色
	Fonts       []string `json:"fonts" yaml:"fonts"`           // 字体列表
}

// MathConfig 数学验证码配置
type MathConfig struct {
	Height      int    `json:"height" yaml:"height"`           // 图片高度
	Width       int    `json:"width" yaml:"width"`             // 图片宽度
	NoiseCount  int    `json:"noiseCount" yaml:"noiseCount"`   // 干扰线数量
	ShowLineOptions int `json:"showLineOptions" yaml:"showLineOptions"` // 显示线条选项
	BgColor     *color.RGBA `json:"bgColor" yaml:"bgColor"` // 背景色
	Fonts       []string `json:"fonts" yaml:"fonts"`           // 字体列表
}

// ChineseConfig 中文验证码配置
type ChineseConfig struct {
	Height      int    `json:"height" yaml:"height"`           // 图片高度
	Width       int    `json:"width" yaml:"width"`             // 图片宽度
	NoiseCount  int    `json:"noiseCount" yaml:"noiseCount"`   // 干扰线数量
	ShowLineOptions int `json:"showLineOptions" yaml:"showLineOptions"` // 显示线条选项
	Length      int    `json:"length" yaml:"length"`           // 验证码长度
	Source      string `json:"source" yaml:"source"`           // 中文字符源
	BgColor     *color.RGBA `json:"bgColor" yaml:"bgColor"` // 背景色
	Fonts       []string `json:"fonts" yaml:"fonts"`           // 字体列表
}

// AudioConfig 音频验证码配置
type AudioConfig struct {
	Length   int    `json:"length" yaml:"length"`     // 验证码长度
	Language string `json:"language" yaml:"language"` // 语言
}

// RandomConfig 随机验证码配置
type RandomConfig struct {
	EnabledTypes []CaptchaType `json:"enabledTypes" yaml:"enabledTypes"` // 启用的验证码类型
	ExcludeAudio bool          `json:"excludeAudio" yaml:"excludeAudio"` // 是否排除音频验证码
}

// StoreConfig 存储配置
type StoreConfig struct {
	Type     StoreType     `json:"type" yaml:"type"`         // 存储类型
	Expire   time.Duration `json:"expire" yaml:"expire"`     // 过期时间
	KeyPrefix string       `json:"keyPrefix" yaml:"keyPrefix"` // Key前缀
	Redis    RedisConfig   `json:"redis" yaml:"redis"`       // Redis配置
}

// StoreType 存储类型
type StoreType string

const (
	StoreTypeMemory StoreType = "memory" // 内存存储
	StoreTypeRedis  StoreType = "redis"  // Redis存储
)

// RedisConfig Redis配置
type RedisConfig struct {
	Addr     string `json:"addr" yaml:"addr"`         // Redis地址
	Password string `json:"password" yaml:"password"` // Redis密码
	DB       int    `json:"db" yaml:"db"`             // Redis数据库
	PoolSize int    `json:"poolSize" yaml:"poolSize"` // 连接池大小
}

// DefaultConfig 返回默认配置
func DefaultConfig() *CaptchaConfig {
	return &CaptchaConfig{
		Type: TypeDigit,
		Digit: DigitConfig{
			Height:   80,
			Width:    240,
			Length:   5,
			MaxSkew:  0.7,
			DotCount: 80,
		},
		String: StringConfig{
			Height:      60,
			Width:       240,
			NoiseCount:  0,
			ShowLineOptions: 0,
			Length:      4,
			Source:      "1234567890qwertyuioplkjhgfdsazxcvbnm",
			BgColor:     &color.RGBA{R: 254, G: 254, B: 254, A: 254},
			Fonts:       []string{"wqy-microhei.ttc"},
		},
		Math: MathConfig{
			Height:      60,
			Width:       240,
			NoiseCount:  0,
			ShowLineOptions: 0,
			BgColor:     &color.RGBA{R: 254, G: 254, B: 254, A: 254},
			Fonts:       []string{"wqy-microhei.ttc"},
		},
		Chinese: ChineseConfig{
			Height:      60,
			Width:       240,
			NoiseCount:  0,
			ShowLineOptions: 0,
			Length:      2,
			Source:      "设想你在处理消费者的音频输出音频可能无论什么都没有任何输出或者它可能是单声道立体声或是环绕立体声的假设我们需要不同类型的音频输出来处理不同的情况",
			BgColor:     &color.RGBA{R: 254, G: 254, B: 254, A: 254},
			Fonts:       []string{"wqy-microhei.ttc"},
		},
		Audio: AudioConfig{
			Length:   4,
			Language: "zh",
		},
		Random: RandomConfig{
			EnabledTypes: []CaptchaType{TypeDigit, TypeString, TypeMath, TypeChinese},
			ExcludeAudio: true,
		},
		Store: StoreConfig{
			Type:      StoreTypeRedis,
			Expire:    5 * time.Minute,
			KeyPrefix: "captcha:",
			Redis: RedisConfig{
				Addr:     "localhost:6379",
				Password: "ynQ28*8g",
				DB:       0,
				PoolSize: 10,
			},
		},
	}
}
