package aws

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"strconv"
	"strings"

	"goprojectskeleton/src/infrastructure/handlers"
)

// LambdaAdapter converts AWS Lambda API Gateway events to handler contexts.
type LambdaAdapter struct{}

// NewLambdaAdapter creates a new Lambda adapter instance.
func NewLambdaAdapter() *LambdaAdapter {
	return &LambdaAdapter{}
}

// APIGatewayV2HTTPRequest represents an API Gateway HTTP API v2 event
type APIGatewayV2HTTPRequest struct {
	Version               string            `json:"version"`
	RouteKey              string            `json:"routeKey"`
	RawPath               string            `json:"rawPath"`
	RawQueryString        string            `json:"rawQueryString"`
	Headers               map[string]string `json:"headers"`
	QueryStringParameters map[string]string `json:"queryStringParameters"`
	PathParameters        map[string]string `json:"pathParameters"`
	RequestContext        struct {
		HTTP struct {
			Method    string `json:"method"`
			Path      string `json:"path"`
			Protocol  string `json:"protocol"`
			SourceIP  string `json:"sourceIp"`
			UserAgent string `json:"userAgent"`
		} `json:"http"`
		RequestID string `json:"requestId"`
	} `json:"requestContext"`
	Body            string `json:"body"`
	IsBase64Encoded bool   `json:"isBase64Encoded"`
}

// APIGatewayV1ProxyRequest represents an API Gateway REST API v1 event
type APIGatewayV1ProxyRequest struct {
	Resource                        string              `json:"resource"`
	Path                            string              `json:"path"`
	HTTPMethod                      string              `json:"httpMethod"`
	Headers                         map[string]string   `json:"headers"`
	MultiValueHeaders               map[string][]string `json:"multiValueHeaders"`
	QueryStringParameters           map[string]string   `json:"queryStringParameters"`
	MultiValueQueryStringParameters map[string][]string `json:"multiValueQueryStringParameters"`
	PathParameters                  map[string]string   `json:"pathParameters"`
	StageVariables                  map[string]string   `json:"stageVariables"`
	RequestContext                  struct {
		RequestID string `json:"requestId"`
		Identity  struct {
			SourceIP string `json:"sourceIp"`
		} `json:"identity"`
	} `json:"requestContext"`
	Body            string `json:"body"`
	IsBase64Encoded bool   `json:"isBase64Encoded"`
}

// LambdaResponse represents the response from a Lambda function
type LambdaResponse struct {
	StatusCode        int                 `json:"statusCode"`
	Headers           map[string]string   `json:"headers"`
	MultiValueHeaders map[string][]string `json:"multiValueHeaders,omitempty"`
	Body              string              `json:"body"`
	IsBase64Encoded   bool                `json:"isBase64Encoded"`
}

// ToHandlerContext converts a Lambda API Gateway event to a handler context.
// This method handles both API Gateway v1 (REST API) and v2 (HTTP API) events.
func (a *LambdaAdapter) ToHandlerContext(
	event interface{},
	responseWriter *LambdaResponseWriter,
	params map[string]string,
) handlers.HandlerContext {
	ctx := context.Background()
	locale := "en-US"
	var body io.ReadCloser
	var query *handlers.Query

	// Try to detect event type and extract information
	switch e := event.(type) {
	case *APIGatewayV2HTTPRequest:
		// API Gateway HTTP API v2
		if acceptLang, ok := e.Headers["accept-language"]; ok {
			locale = acceptLang
		} else if acceptLang, ok := e.Headers["Accept-Language"]; ok {
			locale = acceptLang
		}

		// Parse query parameters
		query = a.parseQueryParamsV2(e.QueryStringParameters, e.RawQueryString)

		// Set body
		if e.Body != "" {
			body = io.NopCloser(bytes.NewBufferString(e.Body))
		}

		// Merge path parameters
		if e.PathParameters != nil {
			if params == nil {
				params = make(map[string]string)
			}
			for k, v := range e.PathParameters {
				params[k] = v
			}
		}

	case *APIGatewayV1ProxyRequest:
		// API Gateway REST API v1
		if acceptLang, ok := e.Headers["Accept-Language"]; ok {
			locale = acceptLang
		} else if acceptLang, ok := e.Headers["accept-language"]; ok {
			locale = acceptLang
		}

		// Parse query parameters
		query = a.parseQueryParamsV1(e.QueryStringParameters, e.MultiValueQueryStringParameters)

		// Set body
		if e.Body != "" {
			body = io.NopCloser(bytes.NewBufferString(e.Body))
		}

		// Merge path parameters
		if e.PathParameters != nil {
			if params == nil {
				params = make(map[string]string)
			}
			for k, v := range e.PathParameters {
				params[k] = v
			}
		}

	case map[string]interface{}:
		// Try to auto-detect event type from map
		if version, ok := e["version"].(string); ok && version == "2.0" {
			// API Gateway v2
			var v2Event APIGatewayV2HTTPRequest
			if eventBytes, err := json.Marshal(event); err == nil {
				if err := json.Unmarshal(eventBytes, &v2Event); err == nil {
					return a.ToHandlerContext(&v2Event, responseWriter, params)
				}
			}
		} else if _, ok := e["httpMethod"].(string); ok {
			// API Gateway v1
			var v1Event APIGatewayV1ProxyRequest
			if eventBytes, err := json.Marshal(event); err == nil {
				if err := json.Unmarshal(eventBytes, &v1Event); err == nil {
					return a.ToHandlerContext(&v1Event, responseWriter, params)
				}
			}
		}
	}

	// Create a mock ResponseWriter for Lambda
	if responseWriter == nil {
		responseWriter = NewLambdaResponseWriter()
	}

	return handlers.NewHandlerContext(
		ctx,
		&locale,
		params,
		&body,
		query,
		responseWriter,
	)
}

// ParsePathParams extracts path parameters from a route pattern and path.
// This is compatible with the Adapter interface but adapted for Lambda.
func (a *LambdaAdapter) ParsePathParams(pattern string, path string) map[string]string {
	params := make(map[string]string)
	patternParts := strings.Split(strings.Trim(pattern, "/"), "/")
	pathParts := strings.Split(strings.Trim(path, "/"), "/")

	if len(patternParts) != len(pathParts) {
		return params
	}

	for i, part := range patternParts {
		if strings.HasPrefix(part, ":") || strings.HasPrefix(part, "{") {
			key := strings.TrimPrefix(part, ":")
			key = strings.TrimPrefix(key, "{")
			key = strings.TrimSuffix(key, "}")
			if i < len(pathParts) {
				params[key] = pathParts[i]
			}
		}
	}

	return params
}

// parseQueryParamsV2 parses query parameters from API Gateway v2 format
func (a *LambdaAdapter) parseQueryParamsV2(queryParams map[string]string, rawQuery string) *handlers.Query {
	if queryParams == nil && rawQuery == "" {
		return nil
	}

	// Parse raw query string if queryParams is empty
	if queryParams == nil && rawQuery != "" {
		queryParams = make(map[string]string)
		parts := strings.Split(rawQuery, "&")
		for _, part := range parts {
			if kv := strings.SplitN(part, "=", 2); len(kv) == 2 {
				queryParams[kv[0]] = kv[1]
			}
		}
	}

	return a.buildQuery(queryParams)
}

// parseQueryParamsV1 parses query parameters from API Gateway v1 format
func (a *LambdaAdapter) parseQueryParamsV1(
	queryParams map[string]string,
	multiValueQueryParams map[string][]string,
) *handlers.Query {
	// Use multi-value if available, otherwise use single-value
	params := make(map[string]string)
	if multiValueQueryParams != nil {
		for k, v := range multiValueQueryParams {
			if len(v) > 0 {
				params[k] = v[0] // Take first value
			}
		}
	} else if queryParams != nil {
		params = queryParams
	} else {
		return nil
	}

	return a.buildQuery(params)
}

// buildQuery builds a Query object from query parameters
func (a *LambdaAdapter) buildQuery(queryParams map[string]string) *handlers.Query {
	var filters []string
	var sorts []string
	page := 0
	pageSize := 0

	if filter, ok := queryParams["filter"]; ok && filter != "" {
		filters = strings.Split(filter, ",")
	}
	if sort, ok := queryParams["sort"]; ok && sort != "" {
		sorts = strings.Split(sort, ",")
	}
	if p, ok := queryParams["page"]; ok && p != "" {
		if parsed, err := strconv.Atoi(p); err == nil {
			page = parsed
		}
	}
	if ps, ok := queryParams["page_size"]; ok && ps != "" {
		if parsed, err := strconv.Atoi(ps); err == nil {
			pageSize = parsed
		}
	}

	if len(filters) > 0 || len(sorts) > 0 || page > 0 || pageSize > 0 {
		return &handlers.Query{
			Filters:  filters,
			Sorts:    sorts,
			Page:     &page,
			PageSize: &pageSize,
		}
	}

	return nil
}

// LambdaResponseWriter is a ResponseWriter implementation for Lambda responses
type LambdaResponseWriter struct {
	StatusCode int
	Headers    map[string]string
	Body       *bytes.Buffer
}

// NewLambdaResponseWriter creates a new LambdaResponseWriter
func NewLambdaResponseWriter() *LambdaResponseWriter {
	return &LambdaResponseWriter{
		StatusCode: http.StatusOK,
		Headers:    make(map[string]string),
		Body:       &bytes.Buffer{},
	}
}

// Header returns the header map
func (w *LambdaResponseWriter) Header() http.Header {
	header := make(http.Header)
	for k, v := range w.Headers {
		header.Set(k, v)
	}
	return header
}

// Write writes data to the response body
func (w *LambdaResponseWriter) Write(b []byte) (int, error) {
	return w.Body.Write(b)
}

// WriteHeader sets the status code
func (w *LambdaResponseWriter) WriteHeader(statusCode int) {
	w.StatusCode = statusCode
}

// ToLambdaResponse converts the ResponseWriter to a LambdaResponse
func (w *LambdaResponseWriter) ToLambdaResponse() LambdaResponse {
	// Set default content type if not set
	if _, ok := w.Headers["Content-Type"]; !ok {
		w.Headers["Content-Type"] = "application/json"
	}

	return LambdaResponse{
		StatusCode: w.StatusCode,
		Headers:    w.Headers,
		Body:       w.Body.String(),
	}
}

// ToHandlerContextFromJSON converts a JSON Lambda event to HandlerContext
// This is a convenience method that handles JSON unmarshaling
func (a *LambdaAdapter) ToHandlerContextFromJSON(
	eventJSON []byte,
	responseWriter *LambdaResponseWriter,
	params map[string]string,
) (handlers.HandlerContext, error) {
	var eventMap map[string]interface{}
	if err := json.Unmarshal(eventJSON, &eventMap); err != nil {
		return handlers.HandlerContext{}, err
	}

	return a.ToHandlerContext(eventMap, responseWriter, params), nil
}

// HandleLambdaEvent is a helper function that processes a Lambda event and returns a Lambda response.
// This function can be used as the main handler in Lambda functions.
// Example usage:
//
//	func handler(ctx context.Context, event map[string]interface{}) (LambdaResponse, error) {
//		return aws.HandleLambdaEvent(event, func(hc handlers.HandlerContext) {
//			handlers.GetHealthCheck(hc)
//		})
//	}
func HandleLambdaEvent(
	event interface{},
	handlerFunc func(handlers.HandlerContext),
) (LambdaResponse, error) {
	adapter := NewLambdaAdapter()
	responseWriter := NewLambdaResponseWriter()
	params := make(map[string]string)

	// Convert event to handler context
	handlerCtx := adapter.ToHandlerContext(event, responseWriter, params)

	// Execute the handler
	handlerFunc(handlerCtx)

	// Convert response writer to Lambda response
	return responseWriter.ToLambdaResponse(), nil
}

// HandleLambdaEventJSON is a convenience function that handles JSON-encoded Lambda events.
// Example usage:
//
//	func handler(ctx context.Context, eventJSON []byte) (LambdaResponse, error) {
//		return aws.HandleLambdaEventJSON(eventJSON, func(hc handlers.HandlerContext) {
//			handlers.GetHealthCheck(hc)
//		})
//	}
func HandleLambdaEventJSON(
	eventJSON []byte,
	handlerFunc func(handlers.HandlerContext),
) (LambdaResponse, error) {
	var eventMap map[string]interface{}
	if err := json.Unmarshal(eventJSON, &eventMap); err != nil {
		// Return error response
		return LambdaResponse{
			StatusCode: http.StatusBadRequest,
			Headers: map[string]string{
				"Content-Type": "application/json",
			},
			Body: `{"error": "Invalid event JSON"}`,
		}, err
	}

	return HandleLambdaEvent(eventMap, handlerFunc)
}
