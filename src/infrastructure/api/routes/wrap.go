package routes

import (
	"gormgoskeleton/src/infrastructure/handlers"

	"github.com/gin-gonic/gin"
)

func wrapHandler(h func(handlers.HandlerContext)) gin.HandlerFunc {
	return func(c *gin.Context) {
		locale := c.GetHeader("X-Locale")

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

		hContext := handlers.NewHandlerContext(c.Request.Context(),
			&locale,
			params,
			&c.Request.Body,
			query,
			c.Writer,
		)

		h(hContext)
	}
}
