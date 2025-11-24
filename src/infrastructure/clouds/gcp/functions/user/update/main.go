package main

import (
	"net/http"

	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
	"gormgoskeleton/src/infrastructure/gcp"
	"gormgoskeleton/src/infrastructure/handlers"
)

func init() {
	gcp.InitializeInfrastructure()
	functions.HTTP("UpdateUser", UpdateUser)
}

// UpdateUser maneja las peticiones de actualización de usuario
func UpdateUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPatch {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Aplicar middleware de autenticación primero
	handler := gcp.AuthMiddleware(func(w http.ResponseWriter, r *http.Request) {
		adapter := gcp.NewHTTPAdapter()
		// Extraer el ID de la ruta
		path := r.URL.Path
		params := adapter.ParsePathParams("/:id", path)

		ctx := adapter.ToHandlerContext(r, w, params)
		handlers.UpdateUser(ctx)
	})
	handler(w, r)
}
