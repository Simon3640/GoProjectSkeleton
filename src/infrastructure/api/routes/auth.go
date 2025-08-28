package routes

import (
	"net/http"

	"gormgoskeleton/src/application/modules/auth"
	dtos "gormgoskeleton/src/application/shared/DTOs"
	"gormgoskeleton/src/application/shared/locales"
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
// @Param        request body dtos.UserCredentials true "User credentials"
// @Success      200 {object} dtos.Token "Tokens generated successfully"
// @Failure      400 {object} map[string]string "Validation error"
// @Router       /api/auth/login [post]
func login(c *gin.Context) {
	var userCredentials dtos.UserCredentials

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
	content, statusCode := api.NewRequestResolver[dtos.Token]().ResolveDTO(c, uc_result, headers)

	c.JSON(statusCode, content)
}

// access-token-refresh
// @Summary      Refresh JWT access token
// @Description  This endpoint allows a user to refresh their JWT access token using a valid refresh
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        request body string true "Refresh token"
// @Success      200 {object} dtos.Token
// @Failure      400 {object} map[string]string "Validation error"
// @Router       /api/auth/refresh [post]
func refreshAccessToken(c *gin.Context) {
	var refreshToken string

	if err := c.ShouldBindJSON(&refreshToken); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	uc_result := auth.NewAuthenticationRefreshUseCase(providers.Logger,
		providers.JWTProviderInstance,
	).Execute(c, locales.EN_US, refreshToken)
	headers := map[api.HTTPHeaderTypeEnum]string{
		api.CONTENT_TYPE: string(api.APPLICATION_JSON),
	}
	content, statusCode := api.NewRequestResolver[dtos.Token]().ResolveDTO(c, uc_result, headers)

	c.JSON(statusCode, content)
}
