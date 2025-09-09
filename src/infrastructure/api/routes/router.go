package routes

import (
	"gormgoskeleton/src/domain/models"
	"gormgoskeleton/src/infrastructure/api/middlewares"
	"gormgoskeleton/src/infrastructure/handlers"

	// "net/http"

	"github.com/gin-gonic/gin"
)

// func wrapHandler(h http.HandlerFunc) gin.HandlerFunc {
// 	return func(c *gin.Context) {

// 		h(c.Writer, c.Request)
// 	}
// }

func Router(r *gin.RouterGroup) {
	r.GET("/health-check", getHealthCheck)

	private := r.Group("/")
	private.Use(middlewares.AuthMiddleware())
	// User routes
	r.POST("/user", wrapHandler(handlers.CreateUser))
	private.GET("/user/:id", wrapHandler(handlers.GetUser))
	private.PATCH("/user/:id", wrapHandler(handlers.UpdateUser))
	private.DELETE("/user/:id", wrapHandler(handlers.DeleteUser))
	private.GET("/user", middlewares.QueryMidleWare[models.User](), wrapHandler(handlers.GetAllUser))
	r.POST("/user-password", wrapHandler(handlers.CreateUserAndPassword))
	r.POST("/user/activate", wrapHandler(handlers.ActivateUser))

	// Password routes
	private.POST("/password", createPassword)
	r.POST("/password/reset-token", createPasswordToken)

	// Auth routes
	r.POST("/auth/login", wrapHandler(handlers.Login))
	r.POST("/auth/refresh", wrapHandler(handlers.RefreshAccessToken))
	r.GET("/auth/password-reset/:identifier", wrapHandler(handlers.RequestPasswordReset))
	r.GET("/auth/login-otp/:otp", wrapHandler(handlers.LoginOTP))

}
