package main

import (
	"fmt"
	"log"

	"github.com/wenpiner/last-admin-common/utils/captcha"
)

func main() {
	// 演示随机验证码功能
	fmt.Println("=== 验证码工具包演示 ===")
	
	// 1. 使用默认配置生成数字验证码
	fmt.Println("\n1. 生成数字验证码:")
	digitDemo()
	
	// 2. 生成随机类型验证码
	fmt.Println("\n2. 生成随机类型验证码:")
	randomDemo()
	
	// 3. 使用全局服务
	fmt.Println("\n3. 使用全局服务:")
	globalServiceDemo()
	
	// 4. 自定义随机配置
	fmt.Println("\n4. 自定义随机配置:")
	customRandomDemo()
}

func digitDemo() {
	config := captcha.DefaultConfig()
	config.Store.Type = captcha.StoreTypeMemory
	config.Type = captcha.TypeDigit

	service, err := captcha.NewService(config)
	if err != nil {
		log.Fatal(err)
	}
	defer service.Close()

	result, err := service.Generate()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("  验证码ID: %s\n", result.ID)
	fmt.Printf("  验证码答案: %s\n", result.Answer)
	fmt.Printf("  Base64数据长度: %d\n", len(result.Base64Blob))

	// 验证
	isValid := service.VerifyAndClear(result.ID, result.Answer)
	fmt.Printf("  验证结果: %t\n", isValid)
}

func randomDemo() {
	config := captcha.DefaultConfig()
	config.Store.Type = captcha.StoreTypeMemory
	config.Type = captcha.TypeRandom

	service, err := captcha.NewService(config)
	if err != nil {
		log.Fatal(err)
	}
	defer service.Close()

	// 生成5个随机验证码
	for i := 1; i <= 5; i++ {
		result, err := service.GenerateRandom()
		if err != nil {
			log.Printf("  生成第%d个验证码失败: %v", i, err)
			continue
		}

		fmt.Printf("  随机验证码%d - ID: %s, 答案: %s\n", i, result.ID, result.Answer)

		// 验证
		isValid := service.VerifyAndClear(result.ID, result.Answer)
		fmt.Printf("  验证结果: %t\n", isValid)
	}
}

func globalServiceDemo() {
	// 初始化全局服务
	config := captcha.DefaultConfig()
	config.Store.Type = captcha.StoreTypeMemory
	config.Type = captcha.TypeRandom

	err := captcha.InitGlobalService(config)
	if err != nil {
		log.Fatal(err)
	}
	defer captcha.CloseGlobalService()

	// 使用全局函数生成不同类型的验证码
	types := []struct {
		name string
		fn   func() (*captcha.CaptchaResult, error)
	}{
		{"数字", captcha.GenerateDigit},
		{"字符串", captcha.GenerateString},
		{"数学", captcha.GenerateMath},
		{"中文", captcha.GenerateChinese},
		{"随机", captcha.GenerateRandom},
	}

	for _, typ := range types {
		result, err := typ.fn()
		if err != nil {
			log.Printf("  生成%s验证码失败: %v", typ.name, err)
			continue
		}

		fmt.Printf("  %s验证码 - ID: %s, 答案: %s\n", typ.name, result.ID, result.Answer)

		// 验证
		isValid := captcha.VerifyAndClear(result.ID, result.Answer)
		fmt.Printf("  验证结果: %t\n", isValid)
	}
}

func customRandomDemo() {
	// 自定义随机配置，只启用数字和字符串验证码
	config := captcha.DefaultConfig()
	config.Store.Type = captcha.StoreTypeMemory
	config.Type = captcha.TypeRandom
	config.Random.EnabledTypes = []captcha.CaptchaType{
		captcha.TypeDigit,
		captcha.TypeString,
	}
	config.Random.ExcludeAudio = true

	service, err := captcha.NewService(config)
	if err != nil {
		log.Fatal(err)
	}
	defer service.Close()

	fmt.Println("  只启用数字和字符串验证码:")
	for i := 1; i <= 5; i++ {
		result, err := service.GenerateRandom()
		if err != nil {
			log.Printf("  生成第%d个验证码失败: %v", i, err)
			continue
		}

		fmt.Printf("  验证码%d - ID: %s, 答案: %s\n", i, result.ID, result.Answer)

		// 验证
		isValid := service.VerifyAndClear(result.ID, result.Answer)
		fmt.Printf("  验证结果: %t\n", isValid)
	}
}
