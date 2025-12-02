package main

import (
	"net/http"

	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
	"gormgoskeleton/src/infrastructure/gcp"
	"gormgoskeleton/src/infrastructure/handlers"
)

func init() {
	gcp.InitializeInfrastructure()
	functions.HTTP("ActivateUser", ActivateUser)
}

// ActivateUser maneja las peticiones de activación de usuario
func ActivateUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	adapter := gcp.NewHTTPAdapter()
	ctx := adapter.ToHandlerContext(r, w, make(map[string]string))
	handlers.ActivateUser(ctx)
}
