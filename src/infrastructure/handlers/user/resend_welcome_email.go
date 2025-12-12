package userhandlers

import (
	"encoding/json"
	"net/http"

	userdtos "github.com/simon3640/goprojectskeleton/src/application/modules/user/dtos"
	userusecases "github.com/simon3640/goprojectskeleton/src/application/modules/user/use_cases"
	database "github.com/simon3640/goprojectskeleton/src/infrastructure/databases/goprojectskeleton"
	userrepositories "github.com/simon3640/goprojectskeleton/src/infrastructure/databases/goprojectskeleton/repositories/user"
	authrepositories "github.com/simon3640/goprojectskeleton/src/infrastructure/databases/goprojectskeleton/repositories/auth"
	handlers "github.com/simon3640/goprojectskeleton/src/infrastructure/handlers/shared"
	"github.com/simon3640/goprojectskeleton/src/infrastructure/providers"
)

// ResendWelcomeEmail Resend welcome email to user
// @Summary Resend welcome email to user
// @Description This endpoint resends the welcome email with a new activation token to the user
// @Tags User
// @Accept json
// @Produce json
// @Param request body dtos.ResendWelcomeEmailRequest true "Email del usuario"
// @Param Accept-Language header string false "Locale for response messages" Enums(en-US, es-ES) default(en-US)
// @Success 200 {object} bool "Correo de bienvenida reenviado"
// @Failure 400 {object} map[string]string "Error de validación"
// @Failure 404 {object} map[string]string "Usuario no encontrado"
// @Router /api/user/resend-welcome-email [post]
func ResendWelcomeEmail(ctx handlers.HandlerContext) {
	var resendRequest userdtos.ResendWelcomeEmailRequest

	if err := json.NewDecoder(*ctx.Body).Decode(&resendRequest); err != nil {
		http.Error(ctx.ResponseWriter, err.Error(), http.StatusBadRequest)
		return
	}

	ucResult := userusecases.NewResendWelcomeEmailUseCase(
		providers.Logger,
		providers.HashProviderInstance,
		userrepositories.NewUserRepository(database.GoProjectSkeletondb.DB, providers.Logger),
		authrepositories.NewOneTimeTokenRepository(database.GoProjectSkeletondb.DB, providers.Logger),
	).Execute(ctx.Context, ctx.Locale, resendRequest)
	headers := map[handlers.HTTPHeaderTypeEnum]string{
		handlers.CONTENT_TYPE: string(handlers.APPLICATION_JSON),
	}
	handlers.NewRequestResolver[bool]().ResolveDTO(ctx.ResponseWriter, ucResult, headers)
}
