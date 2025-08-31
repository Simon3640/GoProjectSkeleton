package routes

import (
	"net/http"

	"gormgoskeleton/src/application/modules/auth"
	auth_pipes "gormgoskeleton/src/application/modules/auth/pipe"
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
// @Success 	 204 {object} nil "OTP login enabled, OTP Sended to user email or phone"
// @Failure      400 {object} map[string]string "Validation error"
// @Router       /api/auth/login [post]
func login(c *gin.Context) {
	var userCredentials dtos.UserCredentials

	if err := c.ShouldBindJSON(&userCredentials); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	password_repository := repositories.NewPasswordRepository(database.DB, providers.Logger)
	userRepository := repositories.NewUserRepository(database.DB, providers.Logger)
	otpRepository := repositories.NewOneTimePasswordRepository(database.DB, providers.Logger)

	uc_result := auth.NewAuthenticateUseCase(providers.Logger,
		password_repository,
		userRepository,
		otpRepository,
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

// password-reset
// @Summary      Request password reset
// @Description  This endpoint allows a user to request a password reset. An email with a
//
//	one-time token will be sent to the user's registered email address.
//
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        identifier path string true "Provided email or phone number"
// @Success      200 {object} map[string]string "Password reset email sent"
// @Failure      400 {object} map[string]string "Validation error"
// @Router       /api/auth/password-reset/{identifier} [get]
func requestPasswordReset(c *gin.Context) {
	identifier := c.Param("identifier")
	if identifier == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Identifier is required"})
		return
	}

	userRepository := repositories.NewUserRepository(database.DB, providers.Logger)
	oneTimeTokenRepository := repositories.NewOneTimeTokenRepository(database.DB, providers.Logger)

	uc_reset_password_token := auth.NewGetResetPasswordTokenUseCase(
		providers.Logger,
		oneTimeTokenRepository,
		userRepository,
		providers.HashProviderInstance,
	)

	uc_reset_password_token_email := auth.NewGetResetPasswordSendEmailUseCase(
		providers.Logger)

	uc_result := auth_pipes.NewGetResetPasswordPipe(
		c.Request.Context(),
		locales.ES_ES,
		uc_reset_password_token,
		uc_reset_password_token_email,
	).Execute(identifier)
	headers := map[api.HTTPHeaderTypeEnum]string{
		api.CONTENT_TYPE: string(api.APPLICATION_JSON),
	}
	content, statusCode := api.NewRequestResolver[bool]().ResolveDTO(c, uc_result, headers)

	c.JSON(statusCode, content)
}

// OTP login
// @Summary      Login with OTP and get JWT tokens
// @Description  This endpoint allows a user to log in with OTP and receive JWT access and
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        otp path string true "One Time Password"
// @Success      200 {object} dtos.Token "Tokens generated successfully"
// @Failure      400 {object} map[string]string "Validation error"
// @Router       /api/auth/login-otp/{otp} [get]
func loginOTP(c *gin.Context) {
	otp := c.Param("otp")

	userRepository := repositories.NewUserRepository(database.DB, providers.Logger)
	otpRepository := repositories.NewOneTimePasswordRepository(database.DB, providers.Logger)

	uc_result := auth.NewAuthenticateOTPUseCase(providers.Logger,
		userRepository,
		otpRepository,
		providers.HashProviderInstance,
		providers.JWTProviderInstance,
	).Execute(c, locales.EN_US, otp)
	headers := map[api.HTTPHeaderTypeEnum]string{
		api.CONTENT_TYPE: string(api.APPLICATION_JSON),
	}
	content, statusCode := api.NewRequestResolver[dtos.Token]().ResolveDTO(c, uc_result, headers)

	c.JSON(statusCode, content)
}
