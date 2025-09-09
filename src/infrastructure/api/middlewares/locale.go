package middlewares

import (
	"gormgoskeleton/src/application/shared/locales"

	"github.com/gin-gonic/gin"
)

func LocaleMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		locale := c.GetHeader("Accept-Language")
		// convert to LocaleTypeEnum
		if locale != string(locales.EN_US) && locale != string(locales.ES_ES) || locale == "" {
			locale = string(locales.EN_US)
		}
		// store locale in context
		c.Set("locale", locale)
		c.Next()
	}
}
