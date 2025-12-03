package aws

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"

	authusecases "goprojectskeleton/src/application/modules/auth/use_cases"
	app_context "goprojectskeleton/src/application/shared/context"
	"goprojectskeleton/src/application/shared/locales"
	"goprojectskeleton/src/application/shared/status"
	database "goprojectskeleton/src/infrastructure/database/goprojectskeleton"
	"goprojectskeleton/src/infrastructure/handlers"
	"goprojectskeleton/src/infrastructure/providers"
	"goprojectskeleton/src/infrastructure/repositories"
)

// AuthMiddleware validates JWT tokens from the Authorization header and injects user context.
// This is for traditional HTTP handlers.
func AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("Authorization")

		locale := r.Header.Get("Accept-Language")
		if locale == "" {
			locale = "en-US"
		}

		ucResult := authusecases.NewAuthUserUseCase(
			providers.Logger,
			repositories.NewUserRepository(database.GoProjectSkeletondb.DB, providers.Logger),
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

// extractTokenFromLambdaEvent extracts the Authorization token from a Lambda event.
// It handles both API Gateway v1 and v2 event formats.
func extractTokenFromLambdaEvent(event interface{}) string {
	switch e := event.(type) {
	case *APIGatewayV2HTTPRequest:
		// API Gateway HTTP API v2 - headers are lowercase
		if auth, ok := e.Headers["authorization"]; ok {
			return auth
		}
		if auth, ok := e.Headers["Authorization"]; ok {
			return auth
		}
	case *APIGatewayV1ProxyRequest:
		// API Gateway REST API v1
		if auth, ok := e.Headers["Authorization"]; ok {
			return auth
		}
		if auth, ok := e.Headers["authorization"]; ok {
			return auth
		}
	case map[string]interface{}:
		// Try to extract from map
		if headers, ok := e["headers"].(map[string]interface{}); ok {
			if auth, ok := headers["Authorization"].(string); ok {
				return auth
			}
			if auth, ok := headers["authorization"].(string); ok {
				return auth
			}
		}
		// Also check for v2 format
		if headers, ok := e["Headers"].(map[string]interface{}); ok {
			if auth, ok := headers["authorization"].(string); ok {
				return auth
			}
			if auth, ok := headers["Authorization"].(string); ok {
				return auth
			}
		}
	}
	return ""
}

// extractLocaleFromLambdaEvent extracts the Accept-Language header from a Lambda event.
func extractLocaleFromLambdaEvent(event interface{}) string {
	locale := "en-US"
	switch e := event.(type) {
	case *APIGatewayV2HTTPRequest:
		if acceptLang, ok := e.Headers["accept-language"]; ok {
			locale = acceptLang
		} else if acceptLang, ok := e.Headers["Accept-Language"]; ok {
			locale = acceptLang
		}
	case *APIGatewayV1ProxyRequest:
		if acceptLang, ok := e.Headers["Accept-Language"]; ok {
			locale = acceptLang
		} else if acceptLang, ok := e.Headers["accept-language"]; ok {
			locale = acceptLang
		}
	case map[string]interface{}:
		if headers, ok := e["headers"].(map[string]interface{}); ok {
			if acceptLang, ok := headers["Accept-Language"].(string); ok {
				locale = acceptLang
			} else if acceptLang, ok := headers["accept-language"].(string); ok {
				locale = acceptLang
			}
		}
	}
	return locale
}

// getStatusMapping returns the HTTP status code mapping for application statuses.
func getStatusMapping() map[status.ApplicationStatusEnum]int {
	return map[status.ApplicationStatusEnum]int{
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
}

// LambdaAuthMiddleware validates JWT tokens from Lambda events and returns an error response if invalid.
// It returns the authenticated user context and a Lambda response (if error) or nil (if success).
func LambdaAuthMiddleware(event interface{}) (context.Context, *LambdaResponse) {
	token := extractTokenFromLambdaEvent(event)
	locale := extractLocaleFromLambdaEvent(event)

	ctx := context.Background()

	ucResult := authusecases.NewAuthUserUseCase(
		providers.Logger,
		repositories.NewUserRepository(database.GoProjectSkeletondb.DB, providers.Logger),
		providers.JWTProviderInstance,
	).Execute(ctx, locales.LocaleTypeEnum(locale), token)

	if ucResult.HasError() {
		statusMapping := getStatusMapping()
		statusCode := 401
		if code, ok := statusMapping[ucResult.StatusCode]; ok {
			statusCode = code
		}

		errorBody, err := json.Marshal(map[string]any{
			"details": ucResult.Error,
		})
		if err != nil {
			providers.Logger.Error("Failed to marshal error response", err)
			errorBody = []byte(`{"details":"Internal server error"}`)
		}

		return nil, &LambdaResponse{
			StatusCode: statusCode,
			Headers: map[string]string{
				"Content-Type": "application/json",
			},
			Body: string(errorBody),
		}
	}

	user := ucResult.GetData()
	ctx = context.WithValue(ctx, app_context.UserKey, *user)

	return ctx, nil
}

// HandleLambdaEventWithAuth is a helper function that processes a Lambda event with authentication.
// It validates the JWT token and injects the user into the context before calling the handler.
// If authentication fails, it returns an error response.
//
// Example usage:
//
//	func handler(ctx context.Context, event map[string]interface{}) (LambdaResponse, error) {
//		return aws.HandleLambdaEventWithAuth(event, func(hc handlers.HandlerContext) {
//			handlers.GetUser(hc)
//		})
//	}
func HandleLambdaEventWithAuth(
	event interface{},
	handlerFunc func(handlers.HandlerContext),
) (LambdaResponse, error) {
	// Validate authentication
	authCtx, authError := LambdaAuthMiddleware(event)
	if authError != nil {
		return *authError, nil
	}

	// Create adapter and response writer
	adapter := NewLambdaAdapter()
	responseWriter := NewLambdaResponseWriter()
	params := make(map[string]string)

	// Convert event to handler context
	handlerCtx := adapter.ToHandlerContext(event, responseWriter, params)

	// Replace context with authenticated context
	localeStr := string(handlerCtx.Locale)
	handlerCtx = handlers.NewHandlerContext(
		authCtx,
		&localeStr,
		handlerCtx.Params,
		handlerCtx.Body,
		handlerCtx.Query,
		handlerCtx.ResponseWriter,
	)

	// Execute the handler
	handlerFunc(handlerCtx)

	// Convert response writer to Lambda response
	return responseWriter.ToLambdaResponse(), nil
}

// HandleLambdaEventWithOptionalAuth processes a Lambda event with optional authentication.
// If a token is present, it validates it and injects the user. If not, it continues without auth.
// This is useful for endpoints that can work with or without authentication.
func HandleLambdaEventWithOptionalAuth(
	event interface{},
	handlerFunc func(handlers.HandlerContext),
) (LambdaResponse, error) {
	token := extractTokenFromLambdaEvent(event)
	locale := extractLocaleFromLambdaEvent(event)
	ctx := context.Background()

	// Only validate if token is present
	if token != "" && strings.TrimSpace(token) != "" {
		ucResult := authusecases.NewAuthUserUseCase(
			providers.Logger,
			repositories.NewUserRepository(database.GoProjectSkeletondb.DB, providers.Logger),
			providers.JWTProviderInstance,
		).Execute(ctx, locales.LocaleTypeEnum(locale), token)

		if !ucResult.HasError() {
			user := ucResult.GetData()
			ctx = context.WithValue(ctx, app_context.UserKey, *user)
		}
		// If error, continue without auth (don't fail the request)
	}

	// Create adapter and response writer
	adapter := NewLambdaAdapter()
	responseWriter := NewLambdaResponseWriter()
	params := make(map[string]string)

	// Convert event to handler context
	handlerCtx := adapter.ToHandlerContext(event, responseWriter, params)

	// Replace context with (possibly authenticated) context
	localeStr := string(handlerCtx.Locale)
	handlerCtx = handlers.NewHandlerContext(
		ctx,
		&localeStr,
		handlerCtx.Params,
		handlerCtx.Body,
		handlerCtx.Query,
		handlerCtx.ResponseWriter,
	)

	// Execute the handler
	handlerFunc(handlerCtx)

	// Convert response writer to Lambda response
	return responseWriter.ToLambdaResponse(), nil
}
