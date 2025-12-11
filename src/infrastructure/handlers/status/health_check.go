package statushandlers

import (
	"time"

	usecases "github.com/simon3640/goprojectskeleton/src/application/modules/status/use_cases"
	"github.com/simon3640/goprojectskeleton/src/domain/models"
	handlers "github.com/simon3640/goprojectskeleton/src/infrastructure/handlers/shared"
	"github.com/simon3640/goprojectskeleton/src/infrastructure/providers"
)

// GetHealthCheck get the status of the API
// @Summary this endpoint get the status of the API
// @Description Get status of the API
// @Schemes models.Status
// @Tags Status Check
// @Accept json
// @Produce json
// @Param Accept-Language header string false "Locale for response messages" Enums(en-US, es-ES) default(en-US)
// @Success 200 {object} models.Status "Status of the API"
// @Router /api/status [get]
func GetHealthCheck(ctx handlers.HandlerContext) {
	ucResult := usecases.NewGetStatusUseCase(
		providers.Logger,
		providers.NewApiStatusProvider(),
	).Execute(ctx.Context, ctx.Locale, time.Now())

	requestResolver := handlers.NewRequestResolver[models.Status]()
	headers := map[handlers.HTTPHeaderTypeEnum]string{
		handlers.CONTENT_TYPE: string(handlers.APPLICATION_JSON),
	}
	requestResolver.ResolveDTO(ctx.ResponseWriter, ucResult, headers)
}
