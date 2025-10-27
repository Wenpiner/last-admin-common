package langctx

import (
	"context"

	"github.com/zeromicro/go-zero/rest/enums"
	"google.golang.org/grpc/metadata"
)

// GetLangFromContext 从上下文中获取语言
func GetLangFromContext(ctx context.Context) (string, bool) {
	if lang, ok := ctx.Value(string(enums.LangKey)).(string); !ok {
		if md, ok := metadata.FromIncomingContext(ctx); !ok {
			return "", false
		} else {
			if data := md.Get(string(enums.LangKey)); len(data) > 0 {
				return data[0], true
			}
		}
	} else {
		return lang, true
	}
	return "", false
}


// WithLangToContext 将语言写入上下文
func WithLangToContext(ctx context.Context, lang string) context.Context {
	return context.WithValue(ctx, string(enums.LangKey), lang)
}
