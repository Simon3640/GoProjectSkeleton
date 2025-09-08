package handlers

import (
	"encoding/json"
	"net/http"

	"gormgoskeleton/src/application/modules/auth"
	auth_pipes "gormgoskeleton/src/application/modules/auth/pipe"
	dtos "gormgoskeleton/src/application/shared/DTOs"
	"gormgoskeleton/src/application/shared/locales"
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
// @Success      200 {object} dtos.Token "Tokens generated successfully"
// @Success 	 204 {object} nil "OTP login enabled, OTP Sended to user email or phone"
// @Failure      400 {object} map[string]string "Validation error"
// @Router       /auth/login [post]
func Login(w http.ResponseWriter, r *http.Request) {
	var userCredentials dtos.UserCredentials
	if err := json.NewDecoder(r.Body).Decode(&userCredentials); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
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
	).Execute(r.Context(), locales.EN_US, userCredentials)
	headers := map[HTTPHeaderTypeEnum]string{
		CONTENT_TYPE: string(APPLICATION_JSON),
	}
	NewRequestResolver[dtos.Token]().ResolveDTO(w, r, uc_result, headers)
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
// @Router       /auth/refresh [post]
func RefreshAccessToken(w http.ResponseWriter, r *http.Request) {
	var refreshToken string
	if err := json.NewDecoder(r.Body).Decode(&refreshToken); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	uc_result := auth.NewAuthenticationRefreshUseCase(providers.Logger,
		providers.JWTProviderInstance,
	).Execute(r.Context(), locales.EN_US, refreshToken)
	headers := map[HTTPHeaderTypeEnum]string{
		CONTENT_TYPE: string(APPLICATION_JSON),
	}
	NewRequestResolver[dtos.Token]().ResolveDTO(w, r, uc_result, headers)
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
// @Router       /auth/password-reset/{identifier} [get]
func RequestPasswordReset(w http.ResponseWriter, r *http.Request) {
	identifier := r.URL.Query().Get("identifier")
	if identifier == "" {
		http.Error(w, "identifier is required", http.StatusBadRequest)
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
		r.Context(),
		locales.ES_ES,
		uc_reset_password_token,
		uc_reset_password_token_email,
	).Execute(identifier)
	headers := map[HTTPHeaderTypeEnum]string{
		CONTENT_TYPE: string(APPLICATION_JSON),
	}

	NewRequestResolver[bool]().ResolveDTO(w, r, uc_result, headers)
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
// @Router       /auth/login-otp/{otp} [get]
func LoginOTP(w http.ResponseWriter, r *http.Request) {
	otp := r.Header.Get("otp")
	if otp == "" {
		http.Error(w, "otp is required", http.StatusBadRequest)
		return
	}

	userRepository := repositories.NewUserRepository(database.DB, providers.Logger)
	otpRepository := repositories.NewOneTimePasswordRepository(database.DB, providers.Logger)

	uc_result := auth.NewAuthenticateOTPUseCase(providers.Logger,
		userRepository,
		otpRepository,
		providers.HashProviderInstance,
		providers.JWTProviderInstance,
	).Execute(r.Context(), locales.EN_US, otp)
	headers := map[HTTPHeaderTypeEnum]string{
		CONTENT_TYPE: string(APPLICATION_JSON),
	}
	NewRequestResolver[dtos.Token]().ResolveDTO(w, r, uc_result, headers)

}
