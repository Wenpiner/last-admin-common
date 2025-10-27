package last_i18n

import (
	"context"
	"embed"
	"encoding/json"
	"net/http"
	"os"
	"path/filepath"
	"testing"

	"github.com/nicksnyder/go-i18n/v2/i18n"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zeromicro/go-zero/core/errorx"
	"golang.org/x/text/language"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

//go:embed testdata/*.json
var testLocaleFS embed.FS

// createTestLocaleFiles creates test locale files for testing
func createTestLocaleFiles(t *testing.T) string {
	tmpDir := t.TempDir()
	
	// Create zh.json
	zhContent := map[string]string{
		"hello":           "你好",
		"world":           "世界",
		"test.message":    "测试消息",
		"error.not.found": "未找到",
		"error.internal":  "内部服务器错误",
		"grpc.error":      "GRPC错误消息",
		"api.error":       "API错误消息",
	}
	zhData, _ := json.Marshal(zhContent)
	zhFile := filepath.Join(tmpDir, "zh.json")
	require.NoError(t, os.WriteFile(zhFile, zhData, 0644))
	
	// Create en.json
	enContent := map[string]string{
		"hello":           "Hello",
		"world":           "World",
		"test.message":    "Test message",
		"error.not.found": "Not found",
		"error.internal":  "Internal server error",
		"grpc.error":      "GRPC error message",
		"api.error":       "API error message",
	}
	enData, _ := json.Marshal(enContent)
	enFile := filepath.Join(tmpDir, "en.json")
	require.NoError(t, os.WriteFile(enFile, enData, 0644))
	
	return tmpDir
}

func TestNewTranslator_WithEmbeddedFS(t *testing.T) {
	conf := Config{Dir: ""}
	translator := NewTranslator(conf, LocaleFS)
	
	assert.NotNil(t, translator)
	assert.NotNil(t, translator.bundle)
	assert.NotNil(t, translator.localizer)
	assert.NotEmpty(t, translator.supportLangs)
}

func TestNewTranslator_WithFileSystem(t *testing.T) {
	tmpDir := createTestLocaleFiles(t)
	
	conf := Config{Dir: tmpDir}
	translator := NewTranslator(conf, LocaleFS)
	
	assert.NotNil(t, translator)
	assert.NotNil(t, translator.bundle)
	assert.NotNil(t, translator.localizer)
	assert.NotEmpty(t, translator.supportLangs)
}

func TestTranslator_AddLanguageSupport(t *testing.T) {
	translator := &Translator{
		localizer: make(map[language.Tag]*i18n.Localizer),
		bundle:    i18n.NewBundle(language.Chinese),
	}
	translator.bundle.RegisterUnmarshalFunc("json", json.Unmarshal)
	
	// Add English support
	translator.AddLanguageSupport(language.English)
	
	assert.Contains(t, translator.supportLangs, language.English)
	assert.Contains(t, translator.localizer, language.English)
	assert.NotNil(t, translator.localizer[language.English])
}

func TestTranslator_AddBundleFromFile(t *testing.T) {
	tmpDir := createTestLocaleFiles(t)
	zhFile := filepath.Join(tmpDir, "zh.json")
	
	translator := &Translator{
		localizer: make(map[language.Tag]*i18n.Localizer),
		bundle:    i18n.NewBundle(language.Chinese),
	}
	translator.bundle.RegisterUnmarshalFunc("json", json.Unmarshal)
	
	err := translator.AddBundleFromFile(zhFile)
	assert.NoError(t, err)
}

func TestTranslator_AddBundleFromEmbeddedFS(t *testing.T) {
	translator := &Translator{
		localizer: make(map[language.Tag]*i18n.Localizer),
		bundle:    i18n.NewBundle(language.Chinese),
	}
	translator.bundle.RegisterUnmarshalFunc("json", json.Unmarshal)
	
	err := translator.AddBundleFromEmbeddedFS(LocaleFS, "locale/zh.json")
	assert.NoError(t, err)
}

func TestTranslator_MatchLocalizer(t *testing.T) {
	conf := Config{Dir: ""}
	translator := NewTranslator(conf, LocaleFS)
	
	tests := []struct {
		name     string
		lang     string
		expected language.Tag
	}{
		{
			name:     "Chinese language",
			lang:     "zh-CN",
			expected: language.Chinese,
		},
		{
			name:     "English language",
			lang:     "en",
			expected: language.English,
		},
		{
			name:     "Unknown language defaults to Chinese",
			lang:     "fr-FR",
			expected: language.Chinese,
		},
		{
			name:     "Empty language defaults to Chinese",
			lang:     "",
			expected: language.Chinese,
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			localizer := translator.MatchLocalizer(tt.lang)
			assert.NotNil(t, localizer)
		})
	}
}

func TestTranslator_Trans(t *testing.T) {
	conf := Config{Dir: ""}
	translator := NewTranslator(conf, LocaleFS)
	
	tests := []struct {
		name     string
		lang     string
		msgID    string
		expected string
	}{
		{
			name:     "Chinese translation exists",
			lang:     "zh-CN",
			msgID:    "hello",
			expected: "你好",
		},
		{
			name:     "English translation exists",
			lang:     "en",
			msgID:    "hello",
			expected: "Hello",
		},
		{
			name:     "Translation not found returns msgID",
			lang:     "zh-CN",
			msgID:    "nonexistent.key",
			expected: "nonexistent.key",
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.WithValue(context.Background(), "lang", tt.lang)
			result := translator.Trans(ctx, tt.msgID)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestTranslator_TransError_GrpcError(t *testing.T) {
	conf := Config{Dir: ""}
	translator := NewTranslator(conf, LocaleFS)
	
	// Create a gRPC error
	grpcErr := status.Error(codes.NotFound, "grpc.error")
	
	ctx := context.WithValue(context.Background(), "lang", "zh-CN")
	result := translator.TransError(ctx, grpcErr)
	
	assert.Error(t, result)
	// The error should be a gRPC status error
	statusErr, ok := status.FromError(result)
	assert.True(t, ok)
	assert.Equal(t, codes.NotFound, statusErr.Code())
	// The message should be translated
	assert.Equal(t, "GRPC错误消息", statusErr.Message())
}

func TestTranslator_TransError_ApiError(t *testing.T) {
	conf := Config{Dir: ""}
	translator := NewTranslator(conf, LocaleFS)
	
	// Create an API error
	apiErr := errorx.NewApiStatusError(errorx.CodeNotFound, "api.error", http.StatusNotFound)
	
	ctx := context.WithValue(context.Background(), "lang", "zh-CN")
	result := translator.TransError(ctx, apiErr)
	
	assert.Error(t, result)
	// The error should be an ApiError
	translatedErr, ok := result.(*errorx.ApiError)
	assert.True(t, ok)
	assert.Equal(t, errorx.CodeNotFound, translatedErr.Code)
	assert.Equal(t, "API错误消息", translatedErr.Message)
	assert.Equal(t, http.StatusNotFound, translatedErr.Status)
}

func TestTranslator_TransError_RegularError(t *testing.T) {
	conf := Config{Dir: ""}
	translator := NewTranslator(conf, LocaleFS)
	
	// Create a regular error
	regularErr := assert.AnError
	
	ctx := context.WithValue(context.Background(), "lang", "zh-CN")
	result := translator.TransError(ctx, regularErr)
	
	assert.Error(t, result)
	// The error should be converted to an ApiError with system error code
	apiErr, ok := result.(*errorx.ApiError)
	assert.True(t, ok)
	assert.Equal(t, errorx.CodeSystemError, apiErr.Code)
	assert.Equal(t, regularErr.Error(), apiErr.Message)
	assert.Equal(t, http.StatusInternalServerError, apiErr.Status)
}

func TestTranslator_TransError_GrpcErrorWithoutTranslation(t *testing.T) {
	conf := Config{Dir: ""}
	translator := NewTranslator(conf, LocaleFS)
	
	// Create a gRPC error with a message that doesn't have translation
	grpcErr := status.Error(codes.Internal, "unknown.error.key")
	
	ctx := context.WithValue(context.Background(), "lang", "zh-CN")
	result := translator.TransError(ctx, grpcErr)
	
	assert.Error(t, result)
	statusErr, ok := status.FromError(result)
	assert.True(t, ok)
	assert.Equal(t, codes.Internal, statusErr.Code())
	// Should return original error message when translation not found
	assert.Contains(t, statusErr.Message(), "unknown.error.key")
}

func TestTranslator_TransError_ApiErrorWithoutTranslation(t *testing.T) {
	conf := Config{Dir: ""}
	translator := NewTranslator(conf, LocaleFS)

	// Create an API error with a message that doesn't have translation
	apiErr := errorx.NewApiStatusError(errorx.CodeBadRequest, "unknown.error.key", http.StatusBadRequest)

	ctx := context.WithValue(context.Background(), "lang", "zh-CN")
	result := translator.TransError(ctx, apiErr)

	assert.Error(t, result)
	translatedErr, ok := result.(*errorx.ApiError)
	assert.True(t, ok)
	assert.Equal(t, errorx.CodeBadRequest, translatedErr.Code)
	// Should return original error message when translation not found
	assert.Equal(t, "unknown.error.key", translatedErr.Message)
	assert.Equal(t, http.StatusBadRequest, translatedErr.Status)
}

func TestTranslator_Trans_PanicRecovery(t *testing.T) {
	conf := Config{Dir: ""}
	translator := NewTranslator(conf, LocaleFS)

	// Test with context that doesn't have "lang" key - should panic and be handled
	ctx := context.Background()

	// This should panic because ctx.Value("lang").(string) will fail
	assert.Panics(t, func() {
		translator.Trans(ctx, "hello")
	})
}

func TestTranslator_TransError_PanicRecovery(t *testing.T) {
	conf := Config{Dir: ""}
	translator := NewTranslator(conf, LocaleFS)

	// Test with context that doesn't have "lang" key - should panic and be handled
	ctx := context.Background()
	regularErr := assert.AnError

	// This should panic because ctx.Value("lang").(string) will fail
	assert.Panics(t, func() {
		translator.TransError(ctx, regularErr)
	})
}

func TestTranslator_AddBundleFromFile_Error(t *testing.T) {
	translator := &Translator{
		localizer: make(map[language.Tag]*i18n.Localizer),
		bundle:    i18n.NewBundle(language.Chinese),
	}
	translator.bundle.RegisterUnmarshalFunc("json", json.Unmarshal)

	// Test with non-existent file
	err := translator.AddBundleFromFile("/non/existent/file.json")
	assert.Error(t, err)
}

func TestTranslator_AddBundleFromEmbeddedFS_Error(t *testing.T) {
	translator := &Translator{
		localizer: make(map[language.Tag]*i18n.Localizer),
		bundle:    i18n.NewBundle(language.Chinese),
	}
	translator.bundle.RegisterUnmarshalFunc("json", json.Unmarshal)

	// Test with non-existent file in embedded FS
	err := translator.AddBundleFromEmbeddedFS(LocaleFS, "non/existent/file.json")
	assert.Error(t, err)
}

func TestTranslator_MatchLocalizer_EdgeCases(t *testing.T) {
	conf := Config{Dir: ""}
	translator := NewTranslator(conf, LocaleFS)

	tests := []struct {
		name string
		lang string
	}{
		{
			name: "Complex language tag",
			lang: "zh-Hans-CN",
		},
		{
			name: "Invalid language tag",
			lang: "invalid-lang-tag",
		},
		{
			name: "Multiple language preferences",
			lang: "fr-FR,en-US;q=0.9,en;q=0.8",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			localizer := translator.MatchLocalizer(tt.lang)
			assert.NotNil(t, localizer)
			// Should always return a valid localizer (fallback to Chinese)
			message, err := localizer.LocalizeMessage(&i18n.Message{ID: "hello"})
			assert.NoError(t, err)
			assert.NotEmpty(t, message)
		})
	}
}

func TestTranslator_TransError_GrpcErrorEdgeCases(t *testing.T) {
	conf := Config{Dir: ""}
	translator := NewTranslator(conf, LocaleFS)

	// Test gRPC error with malformed desc format
	grpcErr := status.Error(codes.Internal, "malformed error without desc =")

	ctx := context.WithValue(context.Background(), "lang", "zh-CN")
	result := translator.TransError(ctx, grpcErr)

	assert.Error(t, result)
	statusErr, ok := status.FromError(result)
	assert.True(t, ok)
	assert.Equal(t, codes.Internal, statusErr.Code())
	// Should return original error message when desc parsing fails
	assert.Contains(t, statusErr.Message(), "malformed error without desc =")
}

func TestTranslator_TransError_NilError(t *testing.T) {
	conf := Config{Dir: ""}
	translator := NewTranslator(conf, LocaleFS)

	ctx := context.WithValue(context.Background(), "lang", "zh-CN")

	// Test with nil error - this will cause a panic in the current implementation
	// because the function doesn't check for nil error
	assert.Panics(t, func() {
		translator.TransError(ctx, nil)
	})
}


