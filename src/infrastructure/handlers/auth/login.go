package authhandlers

import (
	"encoding/json"
	"net/http"

	authdtos "github.com/simon3640/goprojectskeleton/src/application/modules/auth/dtos"
	authusecases "github.com/simon3640/goprojectskeleton/src/application/modules/auth/use_cases"
	database "github.com/simon3640/goprojectskeleton/src/infrastructure/databases/goprojectskeleton"
	authrepositories "github.com/simon3640/goprojectskeleton/src/infrastructure/databases/goprojectskeleton/repositories/auth"
	passwordrepositories "github.com/simon3640/goprojectskeleton/src/infrastructure/databases/goprojectskeleton/repositories/password"
	userrepositories "github.com/simon3640/goprojectskeleton/src/infrastructure/databases/goprojectskeleton/repositories/user"
	handlers "github.com/simon3640/goprojectskeleton/src/infrastructure/handlers/shared"
	"github.com/simon3640/goprojectskeleton/src/infrastructure/providers"
)

// Login login with email and password and get JWT tokens
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
func Login(ctx handlers.HandlerContext) {
	var userCredentials authdtos.UserCredentials
	if err := json.NewDecoder(*ctx.Body).Decode(&userCredentials); err != nil {
		http.Error(ctx.ResponseWriter, err.Error(), http.StatusBadRequest)
		return
	}

	passwordRepository := passwordrepositories.NewPasswordRepository(database.GoProjectSkeletondb.DB, providers.Logger)
	userRepository := userrepositories.NewUserRepository(database.GoProjectSkeletondb.DB, providers.Logger)
	otpRepository := authrepositories.NewOneTimePasswordRepository(database.GoProjectSkeletondb.DB, providers.Logger)

	ucResult := authusecases.NewAuthenticateUseCase(providers.Logger,
		passwordRepository,
		userRepository,
		otpRepository,
		providers.HashProviderInstance,
		providers.JWTProviderInstance,
		providers.CacheProviderInstance,
	).Execute(ctx.Context, ctx.Locale, userCredentials)
	headers := map[handlers.HTTPHeaderTypeEnum]string{
		handlers.CONTENT_TYPE: string(handlers.APPLICATION_JSON),
	}
	handlers.NewRequestResolver[authdtos.Token]().ResolveDTO(ctx.ResponseWriter, ucResult, headers)
}
