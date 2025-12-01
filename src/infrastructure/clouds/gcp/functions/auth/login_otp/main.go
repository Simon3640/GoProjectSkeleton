package main

import (
	"net/http"

	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
	"gormgoskeleton/src/infrastructure/gcp"
	"gormgoskeleton/src/infrastructure/handlers"
)

func init() {
	gcp.InitializeInfrastructure()
	functions.HTTP("LoginOTP", LoginOTP)
}

// LoginOTP maneja las peticiones de login con OTP
func LoginOTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	adapter := gcp.NewHTTPAdapter()
	// Extraer el OTP de la ruta
	path := r.URL.Path
	params := adapter.ParsePathParams("/:otp", path)

	ctx := adapter.ToHandlerContext(r, w, params)
	handlers.LoginOTP(ctx)
}
