package handlers

import (
	"encoding/json"
	"net/http"

	"gormgoskeleton/src/application/modules/auth"
	auth_pipes "gormgoskeleton/src/application/modules/auth/pipe"
	dtos "gormgoskeleton/src/application/shared/DTOs"
	database "gormgoskeleton/src/infrastructure/database/gormgoskeleton"
	"gormgoskeleton/src/infrastructure/providers"
	"gormgoskeleton/src/infrastructure/repositories"
)

// access-token
// @Summary      Login and get JWT tokens
// @Description  This endpoint allows a user to log in and receive JWT access and refresh tokens
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        request body dtos.UserCredentials true "User credentials"
// @Param X-Locale header string false "Locale for response messages" Enums(en-US, es-ES) default(en-US)
// @Success      200 {object} dtos.Token "Tokens generated successfully"
// @Success 	 204 {object} nil "OTP login enabled, OTP Sended to user email or phone"
// @Failure      400 {object} map[string]string "Validation error"
// @Router       /api/auth/login [post]
func Login(ctx HandlerContext) {
	var userCredentials dtos.UserCredentials
	if err := json.NewDecoder(*ctx.Body).Decode(&userCredentials); err != nil {
		http.Error(ctx.ResponseWriter, err.Error(), http.StatusBadRequest)
		return
	}

	password_repository := repositories.NewPasswordRepository(database.DB, providers.Logger)
	userRepository := repositories.NewUserRepository(database.DB, providers.Logger)
	otpRepository := repositories.NewOneTimePasswordRepository(database.DB, providers.Logger)

	ucResult := auth.NewAuthenticateUseCase(providers.Logger,
		password_repository,
		userRepository,
		otpRepository,
		providers.HashProviderInstance,
		providers.JWTProviderInstance,
	).Execute(ctx.c, ctx.Locale, userCredentials)
	headers := map[HTTPHeaderTypeEnum]string{
		CONTENT_TYPE: string(APPLICATION_JSON),
	}
	NewRequestResolver[dtos.Token]().ResolveDTO(ctx.ResponseWriter, ucResult, headers)
}

// access-token-refresh
// @Summary      Refresh JWT access token
// @Description  This endpoint allows a user to refresh their JWT access token using a valid refresh
// @Tags         Auth
// @Accept       json
// @Param X-Locale header string false "Locale for response messages" Enums(en-US, es-ES) default(en-US)
// @Produce      json
// @Param        request body string true "Refresh token"
// @Success      200 {object} dtos.Token
// @Failure      400 {object} map[string]string "Validation error"
// @Router       /api/auth/refresh [post]
func RefreshAccessToken(ctx HandlerContext) {
	var refreshToken string
	if err := json.NewDecoder(*ctx.Body).Decode(&refreshToken); err != nil {
		http.Error(ctx.ResponseWriter, err.Error(), http.StatusBadRequest)
		return
	}

	ucResult := auth.NewAuthenticationRefreshUseCase(providers.Logger,
		providers.JWTProviderInstance,
	).Execute(ctx.c, ctx.Locale, refreshToken)
	headers := map[HTTPHeaderTypeEnum]string{
		CONTENT_TYPE: string(APPLICATION_JSON),
	}
	NewRequestResolver[dtos.Token]().ResolveDTO(ctx.ResponseWriter, ucResult, headers)
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
// @Param X-Locale header string false "Locale for response messages" Enums(en-US, es-ES) default(en-US)
// @Success      200 {object} map[string]string "Password reset email sent"
// @Failure      400 {object} map[string]string "Validation error"
// @Router       /api/auth/password-reset/{identifier} [get]
func RequestPasswordReset(ctx HandlerContext) {
	identifier := ctx.Params["identifier"]
	if identifier == "" {
		http.Error(ctx.ResponseWriter, "identifier is required", http.StatusBadRequest)
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

	ucResult := auth_pipes.NewGetResetPasswordPipe(
		ctx.c,
		ctx.Locale,
		uc_reset_password_token,
		uc_reset_password_token_email,
	).Execute(identifier)
	headers := map[HTTPHeaderTypeEnum]string{
		CONTENT_TYPE: string(APPLICATION_JSON),
	}

	NewRequestResolver[bool]().ResolveDTO(ctx.ResponseWriter, ucResult, headers)
}

// OTP login
// @Summary      Login with OTP and get JWT tokens
// @Description  This endpoint allows a user to log in with OTP and receive JWT access and
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        otp path string true "One Time Password"
// @Param X-Locale header string false "Locale for response messages" Enums(en-US, es-ES) default(en-US)
// @Success      200 {object} dtos.Token "Tokens generated successfully"
// @Failure      400 {object} map[string]string "Validation error"
// @Router       /api/auth/login-otp/{otp} [get]
func LoginOTP(ctx HandlerContext) {
	otp := ctx.Params["otp"]
	if otp == "" {
		http.Error(ctx.ResponseWriter, "otp is required", http.StatusBadRequest)
		return
	}

	userRepository := repositories.NewUserRepository(database.DB, providers.Logger)
	otpRepository := repositories.NewOneTimePasswordRepository(database.DB, providers.Logger)

	ucResult := auth.NewAuthenticateOTPUseCase(providers.Logger,
		userRepository,
		otpRepository,
		providers.HashProviderInstance,
		providers.JWTProviderInstance,
	).Execute(ctx.c, ctx.Locale, otp)
	headers := map[HTTPHeaderTypeEnum]string{
		CONTENT_TYPE: string(APPLICATION_JSON),
	}
	NewRequestResolver[dtos.Token]().ResolveDTO(ctx.ResponseWriter, ucResult, headers)

}
