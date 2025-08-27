package routes

import (
	"net/http"

	usecases_password "gormgoskeleton/src/application/modules/password/use_cases"
	dtos "gormgoskeleton/src/application/shared/DTOs"
	"gormgoskeleton/src/application/shared/locales"
	"gormgoskeleton/src/infrastructure/api"
	database "gormgoskeleton/src/infrastructure/database/gormgoskeleton"
	"gormgoskeleton/src/infrastructure/providers"
	"gormgoskeleton/src/infrastructure/repositories"

	"github.com/gin-gonic/gin"
)

// CreatePassword
// @Summary This endpoint Create a new password
// @Description This endpoint Create a new password
// @Schemes models.PasswordCreateNoHash
// @Tags Password
// @Accept json
// @Produce json
// @Param request body models.PasswordCreateNoHash true "Datos del usuario"
// @Success 201 {object} bool "Usuario creado"
// @Failure 400 {object} map[string]string "Error de validación"
// @Router /api/password [post]
// @Security Bearer
func createPassword(c *gin.Context) {
	var passwordCreate dtos.PasswordCreateNoHash

	if err := c.ShouldBindJSON(&passwordCreate); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	passwordRepository := repositories.NewPasswordRepository(database.DB, providers.Logger)

	uc_result := usecases_password.NewCreatePasswordUseCase(providers.Logger,
		passwordRepository, providers.HashProviderInstance,
	).Execute(c.Request.Context(), locales.EN_US, passwordCreate)
	headers := map[api.HTTPHeaderTypeEnum]string{
		api.CONTENT_TYPE: string(api.APPLICATION_JSON),
	}
	content, statusCode := api.NewRequestResolver[bool]().ResolveDTO(c, uc_result, headers)

	c.JSON(statusCode, content)
}
