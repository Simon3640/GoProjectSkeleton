package handlers

import (
	"time"

	usecases "goprojectskeleton/src/application/modules/status/use_cases"
	"goprojectskeleton/src/domain/models"
	"goprojectskeleton/src/infrastructure/providers"
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
func GetHealthCheck(ctx HandlerContext) {
	ucResult := usecases.NewGetStatusUseCase(
		providers.Logger,
		providers.NewApiStatusProvider(),
	).Execute(ctx.c, ctx.Locale, time.Now())

	requestResolver := NewRequestResolver[models.Status]()
	headers := map[HTTPHeaderTypeEnum]string{
		CONTENT_TYPE: string(APPLICATION_JSON),
	}
	requestResolver.ResolveDTO(ctx.ResponseWriter, ucResult, headers)

}
