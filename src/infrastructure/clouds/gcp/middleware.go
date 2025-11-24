// Package gcp provides Google Cloud Platform specific implementations for the application.
package gcp

import (
	"context"
	"encoding/json"
	"net/http"

	"gormgoskeleton/src/application/modules/auth"
	app_context "gormgoskeleton/src/application/shared/context"
	"gormgoskeleton/src/application/shared/locales"
	"gormgoskeleton/src/application/shared/status"
	database "gormgoskeleton/src/infrastructure/database/gormgoskeleton"
	"gormgoskeleton/src/infrastructure/providers"
	"gormgoskeleton/src/infrastructure/repositories"
)

// AuthMiddleware validates JWT tokens from the Authorization header and injects user context.
func AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("Authorization")

		locale := r.Header.Get("Accept-Language")
		if locale == "" {
			locale = "en-US"
		}

		ucResult := auth.NewAuthUserUseCase(
			providers.Logger,
			repositories.NewUserRepository(database.DB, providers.Logger),
			providers.JWTProviderInstance,
		).Execute(r.Context(), locales.LocaleTypeEnum(locale), token)

		if ucResult.HasError() {
			w.Header().Set("Content-Type", "application/json")
			statusMapping := map[status.ApplicationStatusEnum]int{
				status.Success:                   200,
				status.Updated:                   200,
				status.Created:                   201,
				status.PartialContent:            206,
				status.InvalidInput:              400,
				status.Unauthorized:              401,
				status.NotFound:                  404,
				status.Conflict:                  409,
				status.InternalError:             500,
				status.NotImplemented:            501,
				status.ProviderError:             502,
				status.ChatProviderError:         502,
				status.ProviderEmptyResponse:     502,
				status.ProviderEmptyCacheContext: 502,
			}
			statusCode := 401
			if code, ok := statusMapping[ucResult.StatusCode]; ok {
				statusCode = code
			}
			w.WriteHeader(statusCode)
			if err := json.NewEncoder(w).Encode(map[string]any{
				"details": ucResult.Error,
			}); err != nil {
				providers.Logger.Error("Failed to encode error response", err)
			}
			return
		}

		user := ucResult.GetData()
		ctx := context.WithValue(r.Context(), app_context.UserKey, *user)
		r = r.WithContext(ctx)

		next(w, r)
	}
}
