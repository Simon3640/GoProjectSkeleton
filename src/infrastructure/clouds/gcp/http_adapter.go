// Package gcp provides Google Cloud Platform specific implementations for the application.
package gcp

import (
	"gormgoskeleton/src/adapters"
)

// NewHTTPAdapter creates a new HTTP adapter for GCP.
func NewHTTPAdapter() adapters.Adapter {
	return adapters.NewHTTPAdapter()
}

// GetDefaultAdapter returns the default adapter for GCP.
func GetDefaultAdapter() adapters.Adapter {
	return adapters.GetDefaultAdapter()
}

// CreateAdapter creates an adapter of the specified type for GCP.
func CreateAdapter(adapterType adapters.AdapterType) adapters.Adapter {
	factory := adapters.NewFactory()
	return factory.CreateAdapter(adapterType)
}
