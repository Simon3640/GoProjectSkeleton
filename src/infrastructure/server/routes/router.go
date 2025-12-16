package routes

import (
	"github.com/simon3640/goprojectskeleton/gin/middlewares"
	authhandlers "github.com/simon3640/goprojectskeleton/src/infrastructure/handlers/auth"
	passwordhandlers "github.com/simon3640/goprojectskeleton/src/infrastructure/handlers/password"
	statushandlers "github.com/simon3640/goprojectskeleton/src/infrastructure/handlers/status"
	userhandlers "github.com/simon3640/goprojectskeleton/src/infrastructure/handlers/user"
	"github.com/simon3640/goprojectskeleton/src/infrastructure/observability"

	"github.com/gin-gonic/gin"
)

var observabilityComponents *observability.ObservabilityComponents

// SetObservabilityComponents establece los componentes de observabilidad para las rutas
func SetObservabilityComponents(components *observability.ObservabilityComponents) {
	observabilityComponents = components
}

func Router(r *gin.RouterGroup) {
	// Ruta de status con observabilidad
	if observabilityComponents != nil && observabilityComponents.Tracer != nil {
		r.GET("/status", wrapHandlerWithObservability(
			statushandlers.GetHealthCheck,
			observabilityComponents.Tracer,
			observabilityComponents.Propagator,
		))
	} else {
		r.GET("/status", wrapHandler(statushandlers.GetHealthCheck))
	}

	private := r.Group("/")
	private.Use(middlewares.AuthMiddleware())
	// User routes
	r.POST("/user", wrapHandler(userhandlers.CreateUser))
	private.GET("/user/:id", wrapHandler(userhandlers.GetUser))
	private.PATCH("/user/:id", wrapHandler(userhandlers.UpdateUser))
	private.DELETE("/user/:id", wrapHandler(userhandlers.DeleteUser))
	private.GET("/user", middlewares.QueryMiddleware(), wrapHandler(userhandlers.GetAllUser))
	r.POST("/user-password", wrapHandler(userhandlers.CreateUserAndPassword))
	r.POST("/user/activate", wrapHandler(userhandlers.ActivateUser))
	r.POST("/user/resend-welcome-email", wrapHandler(userhandlers.ResendWelcomeEmail))

	// Password routes
	private.POST("/password", wrapHandler(passwordhandlers.CreatePassword))
	r.POST("/password/reset-token", wrapHandler(passwordhandlers.CreatePasswordToken))

	// Auth routes
	r.POST("/auth/login", wrapHandler(authhandlers.Login))
	r.POST("/auth/refresh", wrapHandler(authhandlers.RefreshAccessToken))
	r.GET("/auth/password-reset/:identifier", wrapHandler(authhandlers.RequestPasswordReset))
	r.GET("/auth/login-otp/:otp", wrapHandler(authhandlers.LoginOTP))

}
