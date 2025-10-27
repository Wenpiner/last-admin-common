package last_i18n

import (
	"context"
	"embed"
	"encoding/json"
	"fmt"
	"io/fs"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/nicksnyder/go-i18n/v2/i18n"
	"github.com/wenpiner/last-admin-common/ctx/langctx"
	"github.com/wenpiner/last-admin-common/utils/errorcode"
	"github.com/wenpiner/last-admin-common/utils/parse"
	"github.com/zeromicro/go-zero/core/errorx"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/rest/enums"
	"golang.org/x/text/language"
	"google.golang.org/grpc/status"
)

//go:embed locale/*.json
var LocaleFS embed.FS

// Translator is a struct storing translating data.
type Translator struct {
	logx.Logger
	bundle       *i18n.Bundle
	localizer    map[language.Tag]*i18n.Localizer
	supportLangs []language.Tag
}

// AddBundleFromEmbeddedFS adds new bundle into translator from embedded file system
func (l *Translator) AddBundleFromEmbeddedFS(file embed.FS, path string) error {
	if _, err := l.bundle.LoadMessageFileFS(file, path); err != nil {
		return err
	}
	return nil
}

// AddBundleFromFile adds new bundle into translator from file path.
func (l *Translator) AddBundleFromFile(path string) error {
	if _, err := l.bundle.LoadMessageFile(path); err != nil {
		return err
	}
	return nil
}

// AddLanguageSupport adds supports for new language
func (l *Translator) AddLanguageSupport(lang language.Tag) {
	l.supportLangs = append(l.supportLangs, lang)
	l.localizer[lang] = i18n.NewLocalizer(l.bundle, lang.String())
}

// Trans used to translate any i18n string.
func (l *Translator) Trans(ctx context.Context, msgId string) string {
	message, err := l.MatchLocalizer(ctx.Value(string(enums.LangKey)).(string)).LocalizeMessage(&i18n.Message{ID: msgId})
	if err != nil {
		return msgId
	}

	if message == "" {
		return msgId
	}

	return message
}

// TransError translates the error message
func (l *Translator) TransError(ctx context.Context, err error) error {
	lang, ok := langctx.GetLangFromContext(ctx)
	if !ok {
		lang = "zh"
		l.Errorw("无法从上下文获取语言", logx.Field("error", err))
	}
	if errorcode.IsGrpcError(err) {
		message, e := l.MatchLocalizer(lang).LocalizeMessage(&i18n.Message{ID: strings.Split(err.Error(), "desc = ")[1]})
		if e != nil || message == "" {
			message = err.Error()
		}
		return status.Error(status.Code(err), message)
	} else if apiErr, ok := err.(*errorx.ApiError); ok {
		message, e := l.MatchLocalizer(lang).LocalizeMessage(&i18n.Message{ID: apiErr.Error()})
		if e != nil {
			message = apiErr.Error()
		}
		return errorx.NewApiStatusError(apiErr.Code, message, apiErr.Status)
	} else {
		return errorx.NewApiStatusError(errorx.CodeSystemError, err.Error(), http.StatusInternalServerError)
	}
}

// MatchLocalizer used to matcher the localizer in map
func (l *Translator) MatchLocalizer(lang string) *i18n.Localizer {
	tags := parse.ParseTags(lang)
	for _, v := range tags {
		if val, ok := l.localizer[v]; ok {
			return val
		}
	}

	return l.localizer[language.Chinese]
}

// NewTranslator returns a translator by I18n Conf.
// If Conf.Dir is empty, it will load paths in embedded FS.
// If Conf.Dir is not empty, it will load paths joined with Dir path.
// e.g. trans = i18n.NewTranslator(c.I18nConf, i18n2.LocaleFS)
func NewTranslator(conf Config, efs embed.FS) *Translator {
	trans := &Translator{
		Logger: logx.WithContext(context.Background()).WithFields(logx.Field("component", "i18n")),
	}
	trans.localizer = make(map[language.Tag]*i18n.Localizer)
	bundle := i18n.NewBundle(language.Chinese)
	bundle.RegisterUnmarshalFunc("json", json.Unmarshal)
	trans.bundle = bundle

	var files []string
	if conf.Dir == "" {
		if err := fs.WalkDir(efs, ".", func(path string, d fs.DirEntry, err error) error {
			if d == nil {
				logx.Must(fmt.Errorf("wrong directory path: %s", conf.Dir))
			}
			if !d.IsDir() {
				files = append(files, path)
			}

			return err
		}); err != nil {
			logx.Must(fmt.Errorf("failed to get any files in dir: %s, error: %v", conf.Dir, err))
		}

		for _, v := range files {
			languageName := strings.TrimSuffix(filepath.Base(v), ".json")
			trans.AddLanguageSupport(parse.ParseTags(languageName)[0])
			err := trans.AddBundleFromEmbeddedFS(efs, v)
			if err != nil {
				logx.Must(fmt.Errorf("failed to load files from %s for i18n, please check the "+
					"configuration, error: %s", v, err.Error()))
			}
		}
	} else {
		if err := filepath.WalkDir(conf.Dir, func(path string, d fs.DirEntry, err error) error {
			if d == nil {
				logx.Must(fmt.Errorf("wrong directory path: %s", conf.Dir))
			}
			if !d.IsDir() {
				files = append(files, path)
			}

			return err
		}); err != nil {
			logx.Must(fmt.Errorf("failed to get any files in dir: %s, error: %v", conf.Dir, err))
		}

		for _, v := range files {
			languageName := strings.TrimSuffix(filepath.Base(v), ".json")
			trans.AddLanguageSupport(parse.ParseTags(languageName)[0])
			err := trans.AddBundleFromFile(v)
			if err != nil {
				logx.Must(fmt.Errorf("failed to load files from %s for i18n, please check the "+
					"configuration, error: %s", filepath.Join(conf.Dir, v), err.Error()))
			}
		}
	}

	return trans
}
