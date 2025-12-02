package handlers

import (
	"time"

	usecases "goprojectskeleton/src/application/modules/status/use_cases"
	"goprojectskeleton/src/domain/models"
	"goprojectskeleton/src/infrastructure/providers"
)

// Health Check
// @Summary
// @Schemes
// @Tags Health Check
// @Accept json
// @Produce json
// @Param Accept-Language header string false "Locale for response messages" Enums(en-US, es-ES) default(en-US)
// @Router /api/health-check [get]
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
