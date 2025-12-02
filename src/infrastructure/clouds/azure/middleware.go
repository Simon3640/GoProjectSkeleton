package azure

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"

	"gormgoskeleton/src/application/modules/auth"
	app_context "gormgoskeleton/src/application/shared/context"
	"gormgoskeleton/src/application/shared/locales"
	"gormgoskeleton/src/application/shared/status"
	database "gormgoskeleton/src/infrastructure/database/gormgoskeleton"
	"gormgoskeleton/src/infrastructure/handlers"
	"gormgoskeleton/src/infrastructure/providers"
	"gormgoskeleton/src/infrastructure/repositories"
)

func AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("Authorization")

		locale := r.Header.Get("Accept-Language")
		if locale == "" {
			locale = "en-US"
		}

		uc_result := auth.NewAuthUserUseCase(
			providers.Logger,
			repositories.NewUserRepository(database.DB, providers.Logger),
			providers.JWTProviderInstance,
		).Execute(r.Context(), locales.LocaleTypeEnum(locale), token)

		if uc_result.HasError() {
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
			if code, ok := statusMapping[uc_result.StatusCode]; ok {
				statusCode = code
			}
			w.WriteHeader(statusCode)
			json.NewEncoder(w).Encode(map[string]any{
				"details": uc_result.Error,
			})
			return
		}

		user := uc_result.GetData()
		ctx := context.WithValue(r.Context(), app_context.UserKey, *user)
		r = r.WithContext(ctx)

		next(w, r)
	}
}

func QueryMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Parse query parameters
		queryParams := r.URL.Query()
		filters := queryParams["filter"]
		sorts := queryParams["sort"]

		page := 0
		if p := queryParams.Get("page"); p != "" {
			if parsed, err := strconv.Atoi(p); err == nil {
				page = parsed
			}
		}

		pageSize := 0
		if ps := queryParams.Get("page_size"); ps != "" {
			if parsed, err := strconv.Atoi(ps); err == nil {
				pageSize = parsed
			}
		}

		// Create query payload
		var query *handlers.Query
		if len(filters) > 0 || len(sorts) > 0 || page > 0 || pageSize > 0 {
			query = &handlers.Query{
				Filters:  filters,
				Sorts:    sorts,
				Page:     &page,
				PageSize: &pageSize,
			}
		}

		// Store query params in context
		ctx := context.WithValue(r.Context(), "queryParams", query)
		r = r.WithContext(ctx)

		next(w, r)
	}
}

func CORSMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, Accept-Language")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Expose-Headers", "*")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		next(w, r)
	}
}
