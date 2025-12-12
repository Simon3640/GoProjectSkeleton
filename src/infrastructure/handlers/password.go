package handlers

import (
	"encoding/json"
	"net/http"

	passworddtos "github.com/simon3640/goprojectskeleton/src/application/modules/password/dtos"
	usecases_password "github.com/simon3640/goprojectskeleton/src/application/modules/password/use_cases"
	database "github.com/simon3640/goprojectskeleton/src/infrastructure/databases/goprojectskeleton"
	repositories "github.com/simon3640/goprojectskeleton/src/infrastructure/databases/goprojectskeleton/repositories"
	"github.com/simon3640/goprojectskeleton/src/infrastructure/providers"
)

// CreatePassword
// @Summary This endpoint Create a new password
// @Description This endpoint Create a new password
// @Schemes dtos.PasswordCreateNoHash
// @Tags Password
// @Accept json
// @Produce json
// @Param Accept-Language header string false "Locale for response messages" Enums(en-US, es-ES) default(en-US)
// @Param request body passworddtos.PasswordCreateNoHash true "Datos del usuario"
// @Success 201 {object} bool "Usuario creado"
// @Failure 400 {object} map[string]string "Error de validación"
// @Router /api/password [post]
// @Security Bearer
func CreatePassword(ctx HandlerContext) {
	var passwordCreate passworddtos.PasswordCreateNoHash

	if err := json.NewDecoder(*ctx.Body).Decode(&passwordCreate); err != nil {
		http.Error(ctx.ResponseWriter, err.Error(), http.StatusBadRequest)
		return
	}

	passwordRepository := repositories.NewPasswordRepository(database.GoProjectSkeletondb.DB, providers.Logger)

	ucResult := usecases_password.NewCreatePasswordUseCase(providers.Logger,
		passwordRepository, providers.HashProviderInstance, false,
	).Execute(ctx.c, ctx.Locale, passwordCreate)
	headers := map[HTTPHeaderTypeEnum]string{
		CONTENT_TYPE: string(APPLICATION_JSON),
	}

	NewRequestResolver[bool]().ResolveDTO(ctx.ResponseWriter, ucResult, headers)
}

// CreatePasswordResetToken
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
func CreatePasswordToken(ctx HandlerContext) {
	var passwordTokenCreate passworddtos.PasswordTokenCreate

	if err := json.NewDecoder(*ctx.Body).Decode(&passwordTokenCreate); err != nil {
		http.Error(ctx.ResponseWriter, err.Error(), http.StatusBadRequest)
		return
	}

	passwordRepository := repositories.NewPasswordRepository(database.GoProjectSkeletondb.DB, providers.Logger)
	oneTimeTokenRepository := repositories.NewOneTimeTokenRepository(database.GoProjectSkeletondb.DB, providers.Logger)

	ucResult := usecases_password.NewCreatePasswordTokenUseCase(providers.Logger,
		passwordRepository,
		providers.HashProviderInstance,
		oneTimeTokenRepository,
	).Execute(ctx.c, ctx.Locale, passwordTokenCreate)

	headers := map[HTTPHeaderTypeEnum]string{
		CONTENT_TYPE: string(APPLICATION_JSON),
	}
	NewRequestResolver[bool]().ResolveDTO(ctx.ResponseWriter, ucResult, headers)
}
