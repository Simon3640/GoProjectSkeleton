package routes

import "github.com/gin-gonic/gin"

func Router(r *gin.RouterGroup) {
	r.GET("/health-check", getHealthCheck)
}
