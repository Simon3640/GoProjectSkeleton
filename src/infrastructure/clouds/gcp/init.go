package gcp

import (
	"gormgoskeleton/src/infrastructure"
)

var initialized bool

// InitializeInfrastructure inicializa la infraestructura compartida
// Esta función debe ser llamada una vez al inicio de cada función serverless
func InitializeInfrastructure() {
	if !initialized {
		infrastructure.Initialize()
		initialized = true
	}
}
