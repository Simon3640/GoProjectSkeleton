// Package auth provides middleware for authentication and authorization in AWS Lambda functions.
package auth

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/simon3640/goprojectskeleton/aws"
	authusecases "github.com/simon3640/goprojectskeleton/src/application/modules/auth/use_cases"
	appcontext "github.com/simon3640/goprojectskeleton/src/application/shared/context"
	"github.com/simon3640/goprojectskeleton/src/application/shared/locales"
	"github.com/simon3640/goprojectskeleton/src/application/shared/observability"
	"github.com/simon3640/goprojectskeleton/src/application/shared/status"
	usecase "github.com/simon3640/goprojectskeleton/src/application/shared/use_case"
	"github.com/simon3640/goprojectskeleton/src/application/shared/workers"
	"github.com/simon3640/goprojectskeleton/src/domain/models"
	database "github.com/simon3640/goprojectskeleton/src/infrastructure/databases/goprojectskeleton"
	userrepositories "github.com/simon3640/goprojectskeleton/src/infrastructure/databases/goprojectskeleton/repositories/user"
	handlers "github.com/simon3640/goprojectskeleton/src/infrastructure/handlers/shared"
	"github.com/simon3640/goprojectskeleton/src/infrastructure/providers"
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

		uc := authusecases.NewAuthUserUseCase(
			userrepositories.NewUserRepository(database.GoProjectSkeletondb.DB, providers.Logger),
			providers.JWTProviderInstance,
		)
		ucResult := usecase.InstrumentUseCase(
			uc,
			&appcontext.AppContext{Context: r.Context()},
			locales.LocaleTypeEnum(locale),
			token,
			observability.GetObservabilityComponents().Tracer,
			observability.GetObservabilityComponents().Metrics,
			observability.GetObservabilityComponents().Clock,
			"auth_user_use_case",
		)

		if ucResult.HasError() {
			w.Header().Set("Content-Type", "application/json")
			statusMapping := getStatusMapping()
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
		ctx := context.WithValue(r.Context(), appcontext.UserKey, *user)
		r = r.WithContext(ctx)

		next(w, r)
	}
}

// extractTokenFromLambdaEvent extracts the Authorization token from a Lambda event.
// It handles both API Gateway v1 and v2 event formats.
func extractTokenFromLambdaEvent(event interface{}) string {
	switch e := event.(type) {
	case *aws.APIGatewayV2HTTPRequest:
		// API Gateway HTTP API v2 - headers are lowercase
		if auth, ok := e.Headers["authorization"]; ok {
			return auth
		}
		if auth, ok := e.Headers["Authorization"]; ok {
			return auth
		}
	case *aws.APIGatewayV1ProxyRequest:
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
	case *aws.APIGatewayV2HTTPRequest:
		if acceptLang, ok := e.Headers["accept-language"]; ok {
			locale = acceptLang
		} else if acceptLang, ok := e.Headers["Accept-Language"]; ok {
			locale = acceptLang
		}
	case *aws.APIGatewayV1ProxyRequest:
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
		status.TooManyRequests:           429,
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
func LambdaAuthMiddleware(event interface{}) (context.Context, *aws.LambdaResponse) {
	token := extractTokenFromLambdaEvent(event)
	locale := extractLocaleFromLambdaEvent(event)

	ctx := context.Background()

	uc := authusecases.NewAuthUserUseCase(
		userrepositories.NewUserRepository(database.GoProjectSkeletondb.DB, providers.Logger),
		providers.JWTProviderInstance,
	)
	ucResult := usecase.InstrumentUseCase(
		uc,
		&appcontext.AppContext{Context: ctx},
		locales.LocaleTypeEnum(locale),
		token,
		observability.GetObservabilityComponents().Tracer,
		observability.GetObservabilityComponents().Metrics,
		observability.GetObservabilityComponents().Clock,
		"auth_user_use_case",
	)

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

		return nil, &aws.LambdaResponse{
			StatusCode: statusCode,
			Headers: map[string]string{
				"Content-Type": "application/json",
			},
			Body: string(errorBody),
		}
	}

	user := ucResult.GetData()
	ctx = context.WithValue(ctx, appcontext.UserKey, *user)

	return ctx, nil
}

// HandleLambdaEventWithAuth is a helper function that processes a Lambda event with authentication.
// It validates the JWT token and injects the user into the context before calling the handler.
// If authentication fails, it returns an error response.
//
// Example usage:
//
//	func handler(ctx context.Context, event map[string]interface{}) (aws.LambdaResponse, error) {
//		return auth.HandleLambdaEventWithAuth(event, func(hc handlers.HandlerContext) {
//			handlers.GetUser(hc)
//		})
//	}
func HandleLambdaEventWithAuth(
	event interface{},
	handlerFunc func(handlers.HandlerContext),
) (aws.LambdaResponse, error) {
	// Validate authentication
	authCtx, authError := LambdaAuthMiddleware(event)
	if authError != nil {
		return *authError, nil
	}

	// Create adapter and response writer
	adapter := aws.NewLambdaAdapter()
	responseWriter := aws.NewLambdaResponseWriter()
	params := make(map[string]string)

	// Convert event to handler context
	handlerCtx := adapter.ToHandlerContext(event, responseWriter, params)

	appContext := appcontext.AppContext{Context: authCtx}
	user := authCtx.Value(appcontext.UserKey)
	if user, ok := user.(models.UserWithRole); ok {
		appContext.AddUserToContext(&user)
	}

	// Replace context with authenticated context
	localeStr := string(handlerCtx.Locale)
	handlerCtx = handlers.NewHandlerContext(
		&appContext,
		&localeStr,
		handlerCtx.Params,
		handlerCtx.Body,
		handlerCtx.Query,
		handlerCtx.ResponseWriter,
	)

	// Execute the handler
	handlerFunc(handlerCtx)

	// Wait for pending background tasks to complete before returning.
	// This is necessary because Lambda freezes execution immediately after
	// returning, which would kill any goroutines running in the background.
	if executor := workers.GetBackgroundExecutor(); executor != nil {
		executor.WaitForPendingTasks()
	}

	// Convert response writer to Lambda response
	return responseWriter.ToLambdaResponse(), nil
}

// HandleLambdaEventWithOptionalAuth processes a Lambda event with optional authentication.
// If a token is present, it validates it and injects the user. If not, it continues without auth.
// This is useful for endpoints that can work with or without authentication.
func HandleLambdaEventWithOptionalAuth(
	event interface{},
	handlerFunc func(handlers.HandlerContext),
) (aws.LambdaResponse, error) {
	token := extractTokenFromLambdaEvent(event)
	locale := extractLocaleFromLambdaEvent(event)
	ctx := context.Background()

	// Only validate if token is present
	if token != "" && strings.TrimSpace(token) != "" {
		uc := authusecases.NewAuthUserUseCase(
			userrepositories.NewUserRepository(database.GoProjectSkeletondb.DB, providers.Logger),
			providers.JWTProviderInstance,
		)
		ucResult := usecase.InstrumentUseCase(
			uc,
			&appcontext.AppContext{Context: ctx},
			locales.LocaleTypeEnum(locale),
			token,
			observability.GetObservabilityComponents().Tracer,
			observability.GetObservabilityComponents().Metrics,
			observability.GetObservabilityComponents().Clock,
			"auth_user_use_case",
		)

		if !ucResult.HasError() {
			user := ucResult.GetData()
			ctx = context.WithValue(ctx, appcontext.UserKey, *user)
		}
		// If error, continue without auth (don't fail the request)
	}

	appContext := appcontext.AppContext{Context: ctx}
	user := ctx.Value(appcontext.UserKey)
	if user, ok := user.(models.UserWithRole); ok {
		appContext.AddUserToContext(&user)
	}

	// Create adapter and response writer
	adapter := aws.NewLambdaAdapter()
	responseWriter := aws.NewLambdaResponseWriter()
	params := make(map[string]string)

	// Convert event to handler context
	handlerCtx := adapter.ToHandlerContext(event, responseWriter, params)

	// Replace context with (possibly authenticated) context
	localeStr := string(handlerCtx.Locale)
	handlerCtx = handlers.NewHandlerContext(
		&appcontext.AppContext{Context: ctx},
		&localeStr,
		handlerCtx.Params,
		handlerCtx.Body,
		handlerCtx.Query,
		handlerCtx.ResponseWriter,
	)

	// Execute the handler
	handlerFunc(handlerCtx)

	// Wait for pending background tasks to complete before returning.
	// This is necessary because Lambda freezes execution immediately after
	// returning, which would kill any goroutines running in the background.
	if executor := workers.GetBackgroundExecutor(); executor != nil {
		executor.WaitForPendingTasks()
	}

	// Convert response writer to Lambda response
	return responseWriter.ToLambdaResponse(), nil
}
