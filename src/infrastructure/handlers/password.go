package handlers

import (
	"encoding/json"
	"net/http"

	usecases_password "gormgoskeleton/src/application/modules/password/use_cases"
	dtos "gormgoskeleton/src/application/shared/DTOs"
	database "gormgoskeleton/src/infrastructure/database/gormgoskeleton"
	"gormgoskeleton/src/infrastructure/providers"
	"gormgoskeleton/src/infrastructure/repositories"
)

// CreatePassword
// @Summary This endpoint Create a new password
// @Description This endpoint Create a new password
// @Schemes dtos.PasswordCreateNoHash
// @Tags Password
// @Accept json
// @Produce json
// @Param request body dtos.PasswordCreateNoHash true "Datos del usuario"
// @Success 201 {object} bool "Usuario creado"
// @Failure 400 {object} map[string]string "Error de validación"
// @Router /api/password [post]
// @Security Bearer
func CreatePassword(ctx HandlerContext) {
	var passwordCreate dtos.PasswordCreateNoHash

	if err := json.NewDecoder(*ctx.Body).Decode(&passwordCreate); err != nil {
		http.Error(ctx.ResponseWriter, err.Error(), http.StatusBadRequest)
		return
	}

	passwordRepository := repositories.NewPasswordRepository(database.DB, providers.Logger)

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
// @Schemes dtos.PasswordTokenCreate
// @Tags Password
// @Accept json
// @Produce json
// @Param request body dtos.PasswordTokenCreate true "Datos del usuario"
// @Success 201 {object} bool "Token creado"
// @Failure 400 {object} map[string]string "Error de validación"
// @Router /api/password/reset-token [post]
func CreatePasswordToken(ctx HandlerContext) {
	var passwordTokenCreate dtos.PasswordTokenCreate

	if err := json.NewDecoder(*ctx.Body).Decode(&passwordTokenCreate); err != nil {
		http.Error(ctx.ResponseWriter, err.Error(), http.StatusBadRequest)
		return
	}

	passwordRepository := repositories.NewPasswordRepository(database.DB, providers.Logger)
	oneTimeTokenRepository := repositories.NewOneTimeTokenRepository(database.DB, providers.Logger)

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
