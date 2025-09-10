package routes

import (
	"gormgoskeleton/src/infrastructure/api/middlewares"
	"gormgoskeleton/src/infrastructure/handlers"

	"github.com/gin-gonic/gin"
)

func Router(r *gin.RouterGroup) {
	r.GET("/health-check", wrapHandler(handlers.GetHealthCheck))

	private := r.Group("/")
	private.Use(middlewares.AuthMiddleware())
	// User routes
	r.POST("/user", wrapHandler(handlers.CreateUser))
	private.GET("/user/:id", wrapHandler(handlers.GetUser))
	private.PATCH("/user/:id", wrapHandler(handlers.UpdateUser))
	private.DELETE("/user/:id", wrapHandler(handlers.DeleteUser))
	private.GET("/user", middlewares.QueryMiddleware(), wrapHandler(handlers.GetAllUser))
	r.POST("/user-password", wrapHandler(handlers.CreateUserAndPassword))
	r.POST("/user/activate", wrapHandler(handlers.ActivateUser))

	// Password routes
	private.POST("/password", wrapHandler(handlers.CreatePassword))
	r.POST("/password/reset-token", wrapHandler(handlers.CreatePasswordToken))

	// Auth routes
	r.POST("/auth/login", wrapHandler(handlers.Login))
	r.POST("/auth/refresh", wrapHandler(handlers.RefreshAccessToken))
	r.GET("/auth/password-reset/:identifier", wrapHandler(handlers.RequestPasswordReset))
	r.GET("/auth/login-otp/:otp", wrapHandler(handlers.LoginOTP))

}
