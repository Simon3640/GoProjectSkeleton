package middlewares

import (
	"goprojectskeleton/src/infrastructure/handlers"
	"strconv"

	"github.com/gin-gonic/gin"
)

func QueryMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Parse query parameters
		filters := c.QueryArray("filter")
		sorts := c.QueryArray("sort")
		page, _ := strconv.Atoi(c.Query("page"))
		pageSize, _ := strconv.Atoi(c.Query("page_size"))

		// Create query payload
		queryParams := handlers.Query{
			Filters:  filters,
			Sorts:    sorts,
			Page:     &page,
			PageSize: &pageSize,
		}

		// Store query params in context
		c.Set("queryParams", queryParams)
		c.Next()
	}
}
