package middlewares

import (
	domain_utils "gormgoskeleton/src/domain/utils"
	"strconv"

	"github.com/gin-gonic/gin"
)

func QueryMidleWare[QueryModel any]() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Parse query parameters
		filters := c.QueryArray("filter")
		sorts := c.QueryArray("sort")
		page, _ := strconv.Atoi(c.Query("page"))
		pageSize, _ := strconv.Atoi(c.Query("page_size"))

		// Create query payload
		queryParams := domain_utils.NewQueryPayloadBuilder[QueryModel](sorts, filters, &page, &pageSize)

		// Store query params in context
		c.Set("queryParams", queryParams)
		c.Next()
	}
}
