// Package adapters provides interfaces and implementations for converting different request formats
// to a common handler context format.
package adapters

import (
	"net/http"

	"github.com/simon3640/goprojectskeleton/src/infrastructure/handlers"
)

// Adapter defines the interface for converting requests to handler contexts.
type Adapter interface {
	ToHandlerContext(r *http.Request, w http.ResponseWriter, params map[string]string) handlers.HandlerContext
	ParsePathParams(pattern string, path string) map[string]string
}

// AdapterType represents the type of adapter to use.
type AdapterType string

const (
	// HTTPAdapterType represents the HTTP adapter type.
	HTTPAdapterType AdapterType = "http"
	// CustomAdapterType represents a custom adapter type.
	CustomAdapterType AdapterType = "custom"
)

// Factory creates adapters based on the specified type.
type Factory struct{}

// NewFactory creates a new adapter factory.
func NewFactory() *Factory {
	return &Factory{}
}

// CreateAdapter creates an adapter of the specified type.
func (f *Factory) CreateAdapter(adapterType AdapterType) Adapter {
	switch adapterType {
	case HTTPAdapterType:
		return NewHTTPAdapter()
	default:
		return NewHTTPAdapter()
	}
}

// GetDefaultAdapter returns the default adapter implementation.
func GetDefaultAdapter() Adapter {
	return NewHTTPAdapter()
}
