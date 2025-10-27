package validator

import (
	"errors"
	"net/http"
	"strings"

	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	enTranslations "github.com/go-playground/validator/v10/translations/en"
	zhTranslations "github.com/go-playground/validator/v10/translations/zh"
	"github.com/zeromicro/go-zero/core/errorx"

	"github.com/wenpiner/last-admin-common/ctx/langctx"
	last_i18n "github.com/wenpiner/last-admin-common/last-i18n"
)

type Validator struct {
	Trans    *last_i18n.Translator
	uni      *ut.UniversalTranslator
	validate *validator.Validate
}

func NewValidator(trans *last_i18n.Translator, uni *ut.UniversalTranslator) *Validator {
	return &Validator{
		Trans:    trans,
		uni:      uni,
		validate: validator.New(),
	}
}

func (v *Validator) Validate(r *http.Request, data any) error {

	lang, ok := langctx.GetLangFromContext(r.Context())
	if !ok {
		lang = r.Header.Get("Accept-Language")
	}

	if lang == "" {
		lang = "zh-CN"
	}
	// 提取主要语言
	lang = strings.Split(lang, ",")[0]
	lang = strings.Split(lang, "-")[0]

	trans, ok := v.uni.GetTranslator(lang)
	if !ok {
		trans, _ = v.uni.GetTranslator("zh")
	}

	switch lang {
	case "en":
		enTranslations.RegisterDefaultTranslations(v.validate, trans)
	case "zh":
		zhTranslations.RegisterDefaultTranslations(v.validate, trans)
	default:
		// 默认使用中文
		zhTranslations.RegisterDefaultTranslations(v.validate, trans)
	}

	err := v.validate.StructCtx(r.Context(), data)
	if err != nil {
		var validationErrors validator.ValidationErrors
		if errors.As(err, &validationErrors) {
			if len(validationErrors) > 0 {
				errorMessage := validationErrors[0].Translate(trans)
				return errorx.NewApiInvalidParamsError(errorMessage)
			}
		}
	}

	return nil
}
