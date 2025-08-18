package routes

import (
	"net/http"

	"gormgoskeleton/src/application/modules/auth"
	"gormgoskeleton/src/application/shared/locales"
	"gormgoskeleton/src/domain/models"
	"gormgoskeleton/src/infrastructure/api"
	database "gormgoskeleton/src/infrastructure/database/gormgoskeleton"
	"gormgoskeleton/src/infrastructure/providers"
	"gormgoskeleton/src/infrastructure/repositories"

	"github.com/gin-gonic/gin"
)

// access-token
// @Summary      Login and get JWT tokens
// @Description  This endpoint allows a user to log in and receive JWT access and refresh tokens
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        request body models.UserCredentials true "User credentials"
// @Success      200 {object} models.Token "Tokens generated successfully"
// @Failure      400 {object} map[string]string "Validation error"
// @Router       /api/auth/login [post]
func login(c *gin.Context) {
	var userCredentials models.UserCredentials

	if err := c.ShouldBindJSON(&userCredentials); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	password_repository := repositories.NewPasswordRepository(database.DB, providers.Logger)

	uc_result := auth.NewAuthenticateUseCase(providers.Logger,
		password_repository,
		providers.HashProviderInstance,
		providers.JWTProviderInstance,
	).Execute(c, locales.EN_US, userCredentials)
	headers := map[api.HTTPHeaderTypeEnum]string{
		api.CONTENT_TYPE: string(api.APPLICATION_JSON),
	}
	content, statusCode := api.NewRequestResolver[models.Token]().ResolveDTO(c, uc_result, headers)

	c.JSON(statusCode, content)
}
