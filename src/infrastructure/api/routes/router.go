package routes

import (
	"gormgoskeleton/src/domain/models"
	"gormgoskeleton/src/infrastructure/api/middlewares"

	"github.com/gin-gonic/gin"
)

func Router(r *gin.RouterGroup) {
	r.GET("/health-check", getHealthCheck)

	private := r.Group("/")
	private.Use(middlewares.AuthMiddleware())
	// User routes
	r.POST("/user", createUser)
	private.GET("/user/:id", getUser)
	private.PATCH("/user/:id", updateUser)
	private.DELETE("/user/:id", deleteUser)
	private.GET("/user", middlewares.QueryMidleWare[models.User](), getAllUser)
	r.POST("/user-password", createUserAndPassword)

	// Password routes
	private.POST("/password", createPassword)

	// Auth routes
	r.POST("/auth/login", login)
	r.POST("/auth/refresh", refreshAccessToken)

}
