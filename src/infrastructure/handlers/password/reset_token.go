package passwordhandlers

import (
	"encoding/json"
	"net/http"

	passworddtos "github.com/simon3640/goprojectskeleton/src/application/modules/password/dtos"
	usecases_password "github.com/simon3640/goprojectskeleton/src/application/modules/password/use_cases"
	database "github.com/simon3640/goprojectskeleton/src/infrastructure/databases/goprojectskeleton"
	authrepositories "github.com/simon3640/goprojectskeleton/src/infrastructure/databases/goprojectskeleton/repositories/auth"
	passwordrepositories "github.com/simon3640/goprojectskeleton/src/infrastructure/databases/goprojectskeleton/repositories/password"
	handlers "github.com/simon3640/goprojectskeleton/src/infrastructure/handlers/shared"
	"github.com/simon3640/goprojectskeleton/src/infrastructure/providers"
)

// CreatePasswordToken create a new password reset token
// @Summary This endpoint Create a new password reset token
// @Description This endpoint Create a new password reset token
// @Schemes passworddtos.PasswordTokenCreate
// @Tags Password
// @Accept json
// @Produce json
// @Param Accept-Language header string false "Locale for response messages" Enums(en-US, es-ES) default(en-US)
// @Param request body passworddtos.PasswordTokenCreate true "Datos del usuario"
// @Success 201 {object} bool "Token creado"
// @Failure 400 {object} map[string]string "Error de validación"
// @Router /api/password/reset-token [post]
func CreatePasswordToken(ctx handlers.HandlerContext) {
	var passwordTokenCreate passworddtos.PasswordTokenCreate

	if err := json.NewDecoder(*ctx.Body).Decode(&passwordTokenCreate); err != nil {
		http.Error(ctx.ResponseWriter, err.Error(), http.StatusBadRequest)
		return
	}

	passwordRepository := passwordrepositories.NewPasswordRepository(database.GoProjectSkeletondb.DB, providers.Logger)
	oneTimeTokenRepository := authrepositories.NewOneTimeTokenRepository(database.GoProjectSkeletondb.DB, providers.Logger)

	ucResult := usecases_password.NewCreatePasswordTokenUseCase(providers.Logger,
		passwordRepository,
		providers.HashProviderInstance,
		oneTimeTokenRepository,
	).Execute(ctx.Context, ctx.Locale, passwordTokenCreate)

	headers := map[handlers.HTTPHeaderTypeEnum]string{
		handlers.CONTENT_TYPE: string(handlers.APPLICATION_JSON),
	}
	handlers.NewRequestResolver[bool]().ResolveDTO(ctx.ResponseWriter, ucResult, headers)
}
