package routes

import (
	app_context "github.com/simon3640/goprojectskeleton/src/application/shared/context"
	usermodels "github.com/simon3640/goprojectskeleton/src/domain/user/models"
	handlers "github.com/simon3640/goprojectskeleton/src/infrastructure/handlers/shared"

	"github.com/gin-gonic/gin"
)

func wrapHandler(h func(handlers.HandlerContext)) gin.HandlerFunc {
	return func(c *gin.Context) {
		locale := c.GetHeader("Accept-Language")

		params := make(map[string]string)
		for _, param := range c.Params {
			params[param.Key] = param.Value
		}

		var query *handlers.Query
		if qp, exists := c.Get("queryParams"); exists {
			if castedQP, ok := qp.(handlers.Query); ok {
				query = &castedQP
			}
		}
		appContext := app_context.AppContext{Context: c.Request.Context()}
		user := c.Request.Context().Value(app_context.UserKey)
		if user, ok := user.(usermodels.UserWithRole); ok {
			appContext.AddUserToContext(&user)
		}
		hContext := handlers.NewHandlerContext(&appContext,
			&locale,
			params,
			&c.Request.Body,
			query,
			c.Writer,
		)

		h(hContext)
	}
}
