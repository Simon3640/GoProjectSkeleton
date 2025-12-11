package authhandlers

import (
	"net/http"

	authdtos "github.com/simon3640/goprojectskeleton/src/application/modules/auth/dtos"
	authusecases "github.com/simon3640/goprojectskeleton/src/application/modules/auth/use_cases"
	database "github.com/simon3640/goprojectskeleton/src/infrastructure/databases/goprojectskeleton"
	userrepositories "github.com/simon3640/goprojectskeleton/src/infrastructure/databases/goprojectskeleton/repositories/user"
	authrepositories "github.com/simon3640/goprojectskeleton/src/infrastructure/databases/goprojectskeleton/repositories/auth"
	handlers "github.com/simon3640/goprojectskeleton/src/infrastructure/handlers/shared"
	"github.com/simon3640/goprojectskeleton/src/infrastructure/providers"
)

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
func LoginOTP(ctx handlers.HandlerContext) {
	otp := ctx.Params["otp"]
	if otp == "" {
		http.Error(ctx.ResponseWriter, "otp is required", http.StatusBadRequest)
		return
	}

	userRepository := userrepositories.NewUserRepository(database.GoProjectSkeletondb.DB, providers.Logger)
	otpRepository := authrepositories.NewOneTimePasswordRepository(database.GoProjectSkeletondb.DB, providers.Logger)

	ucResult := authusecases.NewAuthenticateOTPUseCase(providers.Logger,
		userRepository,
		otpRepository,
		providers.HashProviderInstance,
		providers.JWTProviderInstance,
	).Execute(ctx.Context, ctx.Locale, otp)
	headers := map[handlers.HTTPHeaderTypeEnum]string{
		handlers.CONTENT_TYPE: string(handlers.APPLICATION_JSON),
	}
	handlers.NewRequestResolver[authdtos.Token]().ResolveDTO(ctx.ResponseWriter, ucResult, headers)
}
