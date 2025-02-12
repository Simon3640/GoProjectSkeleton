package routes

import (
	"time"

	"github.com/gin-gonic/gin"

	"gormgoskeleton/src/application/shared/locales"
	usecases "gormgoskeleton/src/application/use_cases"
	"gormgoskeleton/src/domain/models"
	"gormgoskeleton/src/infrastructure/api"
	"gormgoskeleton/src/infrastructure/providers"
)

// Health Check
// @Summary
// @Schemes
// @Tags Health Check
// @Accept json
// @Produce json
// @Router /api/health-check [get]
func getHealthCheck(c *gin.Context) {
	uc_result := usecases.NewGetStatusUseCase(
		providers.Logger,
		providers.NewApiStatusProvider(),
	).Execute(c, locales.EN_US, time.Now())

	requestResolver := api.NewRequestResolver[models.Status]()
	headers := map[api.HTTPHeaderTypeEnum]string{
		api.CONTENT_TYPE: string(api.APPLICATION_JSON),
	}
	content, statusCode := requestResolver.ResolveDTO(c, uc_result, headers)

	c.JSON(statusCode, content)
}
