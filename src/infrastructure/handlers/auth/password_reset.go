package authhandlers

import (
	"net/http"

	authpipes "github.com/simon3640/goprojectskeleton/src/application/modules/auth/pipes"
	authusecases "github.com/simon3640/goprojectskeleton/src/application/modules/auth/use_cases"
	database "github.com/simon3640/goprojectskeleton/src/infrastructure/databases/goprojectskeleton"
	authrepositories "github.com/simon3640/goprojectskeleton/src/infrastructure/databases/goprojectskeleton/repositories/auth"
	userrepositories "github.com/simon3640/goprojectskeleton/src/infrastructure/databases/goprojectskeleton/repositories/user"
	handlers "github.com/simon3640/goprojectskeleton/src/infrastructure/handlers/shared"
	"github.com/simon3640/goprojectskeleton/src/infrastructure/providers"
)

// RequestPasswordReset request a password reset and send an email with a one-time token
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
func RequestPasswordReset(ctx handlers.HandlerContext) {
	identifier := ctx.Params["identifier"]
	if identifier == "" {
		http.Error(ctx.ResponseWriter, "identifier is required", http.StatusBadRequest)
		return
	}

	userRepository := userrepositories.NewUserRepository(database.GoProjectSkeletondb.DB, providers.Logger)
	oneTimeTokenRepository := authrepositories.NewOneTimeTokenRepository(database.GoProjectSkeletondb.DB, providers.Logger)

	ucResetPasswordToken := authusecases.NewGetResetPasswordTokenUseCase(
		providers.Logger,
		oneTimeTokenRepository,
		userRepository,
		providers.HashProviderInstance,
	)

	ucResetPasswordTokenEmail := authusecases.NewGetResetPasswordSendEmailUseCase(
		providers.Logger)

	ucResult := authpipes.NewGetResetPasswordPipe(
		ctx.Context,
		ctx.Locale,
		ucResetPasswordToken,
		ucResetPasswordTokenEmail,
	).Execute(identifier)
	headers := map[handlers.HTTPHeaderTypeEnum]string{
		handlers.CONTENT_TYPE: string(handlers.APPLICATION_JSON),
	}

	handlers.NewRequestResolver[bool]().ResolveDTO(ctx.ResponseWriter, ucResult, headers)
}
