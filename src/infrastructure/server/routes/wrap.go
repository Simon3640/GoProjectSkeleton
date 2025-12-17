package routes

import (
	"fmt"

	"github.com/simon3640/goprojectskeleton/src/application/contracts/observability"
	app_context "github.com/simon3640/goprojectskeleton/src/application/shared/context"
	"github.com/simon3640/goprojectskeleton/src/domain/models"
	handlers "github.com/simon3640/goprojectskeleton/src/infrastructure/handlers/shared"

	"github.com/gin-gonic/gin"
)

func wrapHandler(h func(handlers.HandlerContext)) gin.HandlerFunc {
	return func(c *gin.Context) {
		locale := c.GetHeader("Accept-Language")

		params := make(map[string]string)
		for _, param := range c.Params {
			params[param.Key] = param.Value
		}

		var query *handlers.Query
		if qp, exists := c.Get("queryParams"); exists {
			if castedQP, ok := qp.(handlers.Query); ok {
				query = &castedQP
			}
		}
		appContext := app_context.AppContext{Context: c.Request.Context()}
		user := c.Request.Context().Value(app_context.UserKey)
		if user, ok := user.(models.UserWithRole); ok {
			appContext.AddUserToContext(&user)
		}
		hContext := handlers.NewHandlerContext(&appContext,
			&locale,
			params,
			&c.Request.Body,
			query,
			c.Writer,
		)

		h(hContext)
	}
}

// wrapHandlerWithObservability wraps a handler with observability instrumentation
func wrapHandlerWithObservability(
	h func(handlers.HandlerContext),
	tracer observability.Tracer,
	propagator observability.TracePropagator,
) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Extract trace context from HTTP headers if available
		headers := make(map[string]string)
		for k, v := range c.Request.Header {
			if len(v) > 0 {
				headers[k] = v[0]
			}
		}

		appContext := app_context.AppContext{Context: c.Request.Context()}

		// Try to extract trace context from headers
		if propagator != nil {
			if traceCtx, ok := propagator.Extract(headers); ok && traceCtx.IsValid() {
				appContext.WithTraceContext(traceCtx)
			}
		}

		// Create HTTP request span
		var span observability.Span
		if tracer != nil {
			operation := c.Request.Method + " " + c.FullPath()
			span = tracer.StartSpan(&appContext, "http.request", observability.WithOperation(operation))
			defer span.End()
			span.UpdateAppContext(&appContext)
		}

		locale := c.GetHeader("Accept-Language")

		params := make(map[string]string)
		for _, param := range c.Params {
			params[param.Key] = param.Value
		}

		var query *handlers.Query
		if qp, exists := c.Get("queryParams"); exists {
			if castedQP, ok := qp.(handlers.Query); ok {
				query = &castedQP
			}
		}

		user := c.Request.Context().Value(app_context.UserKey)
		if user, ok := user.(models.UserWithRole); ok {
			appContext.AddUserToContext(&user)
		}

		hContext := handlers.NewHandlerContext(&appContext,
			&locale,
			params,
			&c.Request.Body,
			query,
			c.Writer,
		)

		// Execute handler
		h(hContext)

		// Set span status based on HTTP status code if available
		if span != nil {
			statusCode := c.Writer.Status()
			if statusCode >= 200 && statusCode < 300 {
				span.SetStatus(observability.SpanStatusOK, "")
			} else if statusCode >= 400 {
				span.SetStatus(observability.SpanStatusError, fmt.Sprintf("HTTP %d", statusCode))
			}
		}
	}
}
