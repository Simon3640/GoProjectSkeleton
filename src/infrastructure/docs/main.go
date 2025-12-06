package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"goprojectskeleton/swagger-server/config"
	docs "goprojectskeleton/swagger-server/swagger"

	httpSwagger "github.com/swaggo/http-swagger"
)

// @title Go Project Skeleton API
// @version 1.0
// @description API documentation for GoProjectSkeleton
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8080

// @securityDefinitions.apikey Bearer
// @in header
// @name Authorization
// @Description Bearer token for authentication

func main() {
	// Load configuration
	cfg := config.Load()

	// Configure Swagger information from configuration
	docs.SwaggerInfo.Title = cfg.APITitle
	docs.SwaggerInfo.Version = cfg.APIVersion
	docs.SwaggerInfo.Description = cfg.APIDescription
	docs.SwaggerInfo.Host = cfg.APIHost
	docs.SwaggerInfo.BasePath = cfg.APIBasePath

	// Configure routes
	mux := http.NewServeMux()

	// Servir Swagger UI (httpSwagger maneja automáticamente /swagger.json y /docs/)
	mux.HandleFunc("/docs/", httpSwagger.WrapHandler)
	mux.HandleFunc("/docs", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/docs/", http.StatusMovedPermanently)
	})

	// Root route
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/docs/", http.StatusMovedPermanently)
	})

	// Configure HTTP server
	server := &http.Server{
		Addr:         ":" + cfg.SwaggerPort,
		Handler:      mux,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Configure signal handling for graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	// Start server in goroutine
	go func() {
		log.Printf("Swagger documentation server starting on port %s", cfg.SwaggerPort)
		log.Printf("Swagger UI available at http://localhost:%s/docs/", cfg.SwaggerPort)
		log.Printf("API Host configured: %s", cfg.APIHost)
		log.Printf("API Base Path: %s", cfg.APIBasePath)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Error starting server: %v", err)
		}
	}()

	// Wait for interrupt signal
	<-quit
	log.Println("Shutting down server...")

	// Close server with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exited")
}
