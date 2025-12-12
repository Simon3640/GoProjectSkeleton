package handlers

import (
	"encoding/json"
	"net/http"

	authdtos "github.com/simon3640/goprojectskeleton/src/application/modules/auth/dtos"
	authpipes "github.com/simon3640/goprojectskeleton/src/application/modules/auth/pipes"
	authusecases "github.com/simon3640/goprojectskeleton/src/application/modules/auth/use_cases"
	database "github.com/simon3640/goprojectskeleton/src/infrastructure/databases/goprojectskeleton"
	repositories "github.com/simon3640/goprojectskeleton/src/infrastructure/databases/goprojectskeleton/repositories"
	"github.com/simon3640/goprojectskeleton/src/infrastructure/providers"
)

// access-token
// @Summary      Login and get JWT tokens
// @Description  This endpoint allows a user to log in and receive JWT access and refresh tokens
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        request body authdtos.UserCredentials true "User credentials"
// @Param Accept-Language header string false "Locale for response messages" Enums(en-US, es-ES) default(en-US)
// @Success      200 {object} authdtos.Token "Tokens generated successfully"
// @Success 	 204 {object} nil "OTP login enabled, OTP Sended to user email or phone"
// @Failure      400 {object} map[string]string "Validation error"
// @Router       /api/auth/login [post]
func Login(ctx HandlerContext) {
	var userCredentials authdtos.UserCredentials
	if err := json.NewDecoder(*ctx.Body).Decode(&userCredentials); err != nil {
		http.Error(ctx.ResponseWriter, err.Error(), http.StatusBadRequest)
		return
	}

	passwordRepository := repositories.NewPasswordRepository(database.GoProjectSkeletondb.DB, providers.Logger)
	userRepository := repositories.NewUserRepository(database.GoProjectSkeletondb.DB, providers.Logger)
	otpRepository := repositories.NewOneTimePasswordRepository(database.GoProjectSkeletondb.DB, providers.Logger)

	ucResult := authusecases.NewAuthenticateUseCase(providers.Logger,
		passwordRepository,
		userRepository,
		otpRepository,
		providers.HashProviderInstance,
		providers.JWTProviderInstance,
		providers.CacheProviderInstance,
	).Execute(ctx.c, ctx.Locale, userCredentials)
	headers := map[HTTPHeaderTypeEnum]string{
		CONTENT_TYPE: string(APPLICATION_JSON),
	}
	NewRequestResolver[authdtos.Token]().ResolveDTO(ctx.ResponseWriter, ucResult, headers)
}

// access-token-refresh
// @Summary      Refresh JWT access token
// @Description  This endpoint allows a user to refresh their JWT access token using a valid refresh
// @Tags         Auth
// @Accept       json
// @Param Accept-Language header string false "Locale for response messages" Enums(en-US, es-ES) default(en-US)
// @Produce      json
// @Param        request body string true "Refresh token"
// @Success      200 {object} authdtos.Token
// @Failure      400 {object} map[string]string "Validation error"
// @Router       /api/auth/refresh [post]
func RefreshAccessToken(ctx HandlerContext) {
	var refreshToken string
	if err := json.NewDecoder(*ctx.Body).Decode(&refreshToken); err != nil {
		http.Error(ctx.ResponseWriter, err.Error(), http.StatusBadRequest)
		return
	}

	ucResult := authusecases.NewAuthenticationRefreshUseCase(providers.Logger,
		providers.JWTProviderInstance,
	).Execute(ctx.c, ctx.Locale, refreshToken)
	headers := map[HTTPHeaderTypeEnum]string{
		CONTENT_TYPE: string(APPLICATION_JSON),
	}
	NewRequestResolver[authdtos.Token]().ResolveDTO(ctx.ResponseWriter, ucResult, headers)
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
// @Param Accept-Language header string false "Locale for response messages" Enums(en-US, es-ES) default(en-US)
// @Success      200 {object} map[string]string "Password reset email sent"
// @Failure      400 {object} map[string]string "Validation error"
// @Router       /api/auth/password-reset/{identifier} [get]
func RequestPasswordReset(ctx HandlerContext) {
	identifier := ctx.Params["identifier"]
	if identifier == "" {
		http.Error(ctx.ResponseWriter, "identifier is required", http.StatusBadRequest)
		return
	}

	userRepository := repositories.NewUserRepository(database.GoProjectSkeletondb.DB, providers.Logger)
	oneTimeTokenRepository := repositories.NewOneTimeTokenRepository(database.GoProjectSkeletondb.DB, providers.Logger)

	ucResetPasswordToken := authusecases.NewGetResetPasswordTokenUseCase(
		providers.Logger,
		oneTimeTokenRepository,
		userRepository,
		providers.HashProviderInstance,
	)

	ucResetPasswordTokenEmail := authusecases.NewGetResetPasswordSendEmailUseCase(
		providers.Logger)

	ucResult := authpipes.NewGetResetPasswordPipe(
		ctx.c,
		ctx.Locale,
		ucResetPasswordToken,
		ucResetPasswordTokenEmail,
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
// @Param Accept-Language header string false "Locale for response messages" Enums(en-US, es-ES) default(en-US)
// @Success      200 {object} authdtos.Token "Tokens generated successfully"
// @Failure      400 {object} map[string]string "Validation error"
// @Router       /api/auth/login-otp/{otp} [get]
func LoginOTP(ctx HandlerContext) {
	otp := ctx.Params["otp"]
	if otp == "" {
		http.Error(ctx.ResponseWriter, "otp is required", http.StatusBadRequest)
		return
	}

	userRepository := repositories.NewUserRepository(database.GoProjectSkeletondb.DB, providers.Logger)
	otpRepository := repositories.NewOneTimePasswordRepository(database.GoProjectSkeletondb.DB, providers.Logger)

	ucResult := authusecases.NewAuthenticateOTPUseCase(providers.Logger,
		userRepository,
		otpRepository,
		providers.HashProviderInstance,
		providers.JWTProviderInstance,
	).Execute(ctx.c, ctx.Locale, otp)
	headers := map[HTTPHeaderTypeEnum]string{
		CONTENT_TYPE: string(APPLICATION_JSON),
	}
	NewRequestResolver[authdtos.Token]().ResolveDTO(ctx.ResponseWriter, ucResult, headers)

}
