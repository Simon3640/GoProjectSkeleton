package routes

import "github.com/gin-gonic/gin"

func Router(r *gin.RouterGroup) {
	r.GET("/health-check", getHealthCheck)
	r.POST("/user", createUser)
	r.GET("/user/:id", getUser)
	r.PATCH("/user/:id", updateUser)
	r.DELETE("/user/:id", deleteUser)
	r.GET("/user", getAllUser)
}
