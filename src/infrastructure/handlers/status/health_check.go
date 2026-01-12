// Package statushandlers contains the handlers for the status module
package statushandlers

import (
	"time"

	usecases "github.com/simon3640/goprojectskeleton/src/application/modules/status/use_cases"
	"github.com/simon3640/goprojectskeleton/src/application/shared/observability"
	usecase "github.com/simon3640/goprojectskeleton/src/application/shared/use_case"
	statusmodels "github.com/simon3640/goprojectskeleton/src/domain/status/models"
	handlers "github.com/simon3640/goprojectskeleton/src/infrastructure/handlers/shared"
	"github.com/simon3640/goprojectskeleton/src/infrastructure/providers"
)

// GetHealthCheck get the status of the API
// @Summary this endpoint get the status of the API
// @Description Get status of the API
// @Schemes statusmodels.Status
// @Tags Status Check
// @Accept json
// @Produce json
// @Param Accept-Language header string false "Locale for response messages" Enums(en-US, es-ES) default(en-US)
// @Success 200 {object} statusmodels.Status "Status of the API"
// @Router /api/status [get]
func GetHealthCheck(ctx handlers.HandlerContext) {
	uc := usecases.NewGetStatusUseCase(
		providers.NewApiStatusProvider(),
	)
	ucResult := usecase.InstrumentUseCase(
		uc,
		ctx.Context,
		ctx.Locale,
		time.Now(),
		observability.GetObservabilityComponents().Tracer,
		observability.GetObservabilityComponents().Metrics,
		observability.GetObservabilityComponents().Clock,
		"get_status_use_case",
	)

	requestResolver := handlers.NewRequestResolver[statusmodels.Status]()
	headers := map[handlers.HTTPHeaderTypeEnum]string{
		handlers.CONTENT_TYPE: string(handlers.APPLICATION_JSON),
	}
	requestResolver.ResolveDTO(ctx.ResponseWriter, ucResult, headers)
}
