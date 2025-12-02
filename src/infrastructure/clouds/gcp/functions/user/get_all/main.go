package main

import (
	"net/http"

	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
	"gormgoskeleton/src/infrastructure/gcp"
	"gormgoskeleton/src/infrastructure/handlers"
)

func init() {
	gcp.InitializeInfrastructure()
	functions.HTTP("GetAllUser", GetAllUser)
}

// GetAllUser maneja las peticiones de obtener todos los usuarios
func GetAllUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Aplicar middleware de autenticación primero
	handler := gcp.AuthMiddleware(func(w http.ResponseWriter, r *http.Request) {
		adapter := gcp.NewHTTPAdapter()
		ctx := adapter.ToHandlerContext(r, w, make(map[string]string))
		handlers.GetAllUser(ctx)
	})
	handler(w, r)
}
