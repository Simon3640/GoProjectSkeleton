package main

import (
	"net/http"

	"gormgoskeleton/src/infrastructure/gcp"
	"gormgoskeleton/src/infrastructure/handlers"

	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
)

func init() {
	gcp.InitializeInfrastructure()
	functions.HTTP("Login", Login)
}

func Login(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	adapter := gcp.NewHTTPAdapter()
	ctx := adapter.ToHandlerContext(r, w, make(map[string]string))
	handlers.Login(ctx)
}
