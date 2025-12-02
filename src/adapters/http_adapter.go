package adapters

import (
	"net/http"
	"strconv"
	"strings"

	"goprojectskeleton/src/infrastructure/handlers"
)

// HTTPAdapter converts HTTP requests to handler contexts.
type HTTPAdapter struct{}

// NewHTTPAdapter creates a new HTTP adapter instance.
func NewHTTPAdapter() *HTTPAdapter {
	return &HTTPAdapter{}
}

// ToHandlerContext converts an HTTP request to a handler context.
func (a *HTTPAdapter) ToHandlerContext(
	r *http.Request,
	w http.ResponseWriter,
	params map[string]string,
) handlers.HandlerContext {
	locale := r.Header.Get("Accept-Language")
	if locale == "" {
		locale = "en-US"
	}

	var query *handlers.Query
	// Check if query params are already in context (set by QueryMiddleware)
	if qp := r.Context().Value("queryParams"); qp != nil {
		if castedQP, ok := qp.(*handlers.Query); ok {
			query = castedQP
		}
	}
	// If not in context, parse from URL
	if query == nil && r.URL.Query() != nil {
		query = parseQueryParams(r.URL.Query())
	}

	body := r.Body

	return handlers.NewHandlerContext(
		r.Context(),
		&locale,
		params,
		&body,
		query,
		w,
	)
}

// ParsePathParams extracts path parameters from a URL pattern and path.
func (a *HTTPAdapter) ParsePathParams(pattern string, path string) map[string]string {
	params := make(map[string]string)
	patternParts := strings.Split(strings.Trim(pattern, "/"), "/")
	pathParts := strings.Split(strings.Trim(path, "/"), "/")

	if len(patternParts) != len(pathParts) {
		return params
	}

	for i, part := range patternParts {
		if strings.HasPrefix(part, ":") {
			key := strings.TrimPrefix(part, ":")
			if i < len(pathParts) {
				params[key] = pathParts[i]
			}
		}
	}

	return params
}

func parseQueryParams(queryParams map[string][]string) *handlers.Query {
	filters := queryParams["filter"]
	sorts := queryParams["sort"]
	page := 0
	pageSize := 0

	if p := queryParams["page"]; len(p) > 0 && p[0] != "" {
		if parsed, err := strconv.Atoi(p[0]); err == nil {
			page = parsed
		}
	}

	if ps := queryParams["page_size"]; len(ps) > 0 && ps[0] != "" {
		if parsed, err := strconv.Atoi(ps[0]); err == nil {
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
