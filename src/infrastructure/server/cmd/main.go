package main

import (
	routes "goprojectskeleton/gin/routes"
	"goprojectskeleton/src/application/shared/settings"
	"goprojectskeleton/src/infrastructure"
	providers "goprojectskeleton/src/infrastructure/providers"

	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/graceful"
	"github.com/gin-gonic/gin"
)

func main() {
	providers.Logger.Info("Initializing infraestructure...")
	infrastructure.Initialize()
	providers.Logger.Info("Infraestructure initialized")
	providers.Logger.Info("Starting server...")
	providers.Logger.Info("App Name: " + settings.AppSettingsInstance.AppName)
	providers.Logger.Info("Port: " + settings.AppSettingsInstance.AppPort)
	// ctx := context.Background()

	app := buildGinApp()
	defer app.Close()
	loadGinApp(app)
	if err := app.Run("0.0.0.0:" + settings.AppSettingsInstance.AppPort); err != nil {
		providers.Logger.Panic("Error running server", err)
	}

}

func loadGinApp(app *graceful.Graceful) {
	app.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"*"},
		AllowHeaders:     []string{"*"},
		ExposeHeaders:    []string{"content-disposition", " content-description"},
		AllowCredentials: true,
	}))
	app.NoRoute(func(c *gin.Context) {
		c.JSON(404, gin.H{"message": "Not Found"})
	})
	app.NoMethod(func(c *gin.Context) {
		c.JSON(405, gin.H{"message": "Method Not Allowed"})
	})
	routes.Router(app.Group("/api"))
}

func buildGinApp() *graceful.Graceful {
	gracefulApp, err := graceful.Default()
	if err != nil {
		providers.Logger.Panic("Error creating graceful app", err)
	}
	gracefulApp.Use(
		gin.Recovery())
	return gracefulApp
}
