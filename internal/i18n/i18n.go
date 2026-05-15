package i18n

import (
	"embed"
	"net/http"

	"github.com/BurntSushi/toml"
	"github.com/gin-gonic/gin"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
)

//go:embed locales/*.toml
var localeFS embed.FS

const contextKey = "localizer"

var bundle *i18n.Bundle

func init() {
	bundle = i18n.NewBundle(language.Chinese)
	bundle.RegisterUnmarshalFunc("toml", toml.Unmarshal)
	if _, err := bundle.LoadMessageFileFS(localeFS, "locales/active.zh-CN.toml"); err != nil {
		panic(err)
	}
	if _, err := bundle.LoadMessageFileFS(localeFS, "locales/active.en-US.toml"); err != nil {
		panic(err)
	}
}

// Bundle returns the global i18n bundle.
func Bundle() *i18n.Bundle {
	return bundle
}

// Middleware extracts the Accept-Language header and stores a Localizer in the Gin context.
// Falls back to zh-CN if no header is present.
func Middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		accept := c.GetHeader("Accept-Language")
		if accept == "" {
			accept = "zh-CN"
		}
		localizer := i18n.NewLocalizer(bundle, accept)
		c.Set(contextKey, localizer)
		c.Next()
	}
}

// Localize translates a message ID using the Localizer stored in the Gin context.
func Localize(c *gin.Context, msgID string) string {
	localizer, ok := c.Value(contextKey).(*i18n.Localizer)
	if !ok {
		// Fallback localizer using zh-CN
		localizer = i18n.NewLocalizer(bundle, "zh-CN")
	}
	msg, err := localizer.Localize(&i18n.LocalizeConfig{MessageID: msgID})
	if err != nil {
		return msgID
	}
	return msg
}

// NoRoute404 returns a handler for unmatched routes.
func NoRoute404(c *gin.Context) {
	c.AbortWithStatusJSON(http.StatusNotFound, gin.H{
		"code": 404,
		"msg":  Localize(c, MsgEndpointNotFound),
		"data": nil,
	})
}
