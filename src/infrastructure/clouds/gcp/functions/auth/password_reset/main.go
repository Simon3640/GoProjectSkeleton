package main

import (
	"net/http"

	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
	"gormgoskeleton/src/infrastructure/gcp"
	"gormgoskeleton/src/infrastructure/handlers"
)

func init() {
	gcp.InitializeInfrastructure()
	functions.HTTP("RequestPasswordReset", RequestPasswordReset)
}

// RequestPasswordReset maneja las peticiones de reset de contraseña
func RequestPasswordReset(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	adapter := gcp.NewHTTPAdapter()
	// Extraer el identificador de la ruta
	path := r.URL.Path
	params := adapter.ParsePathParams("/:identifier", path)

	ctx := adapter.ToHandlerContext(r, w, params)
	handlers.RequestPasswordReset(ctx)
}
