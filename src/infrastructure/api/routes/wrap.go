package routes

import (
	"gormgoskeleton/src/application/shared/locales"
	"gormgoskeleton/src/infrastructure/handlers"
	"net/http"

	"github.com/gin-gonic/gin"
)

func wrapHandler[QueryModel any](h http.HandlerFunc) gin.HandlerFunc {
	return func(c *gin.Context) {
		locale := c.GetHeader("Accept-Language")
		if locale == "" {
			locale = "en-US"
		}
		localeEnum := locales.LocaleTypeEnum(locale)

		hContext := handlers.NewHandlerContext(c.Request.Context(), localeEnum, )
}