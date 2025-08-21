package routes

import "github.com/gin-gonic/gin"

func Router(r *gin.RouterGroup) {
	r.GET("/health-check", getHealthCheck)

	// User routes
	r.POST("/user", createUser)
	r.GET("/user/:id", getUser)
	r.PATCH("/user/:id", updateUser)
	r.DELETE("/user/:id", deleteUser)
	r.GET("/user", getAllUser)
	r.POST("/user-password", createUserAndPassword)

	// Password routes
	r.POST("/password", createPassword)

	// Auth routes
	r.POST("/auth/login", login)
	r.POST("/auth/refresh", refreshAccessToken)

}
