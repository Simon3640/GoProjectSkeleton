package passwordhandlers

import (
	"encoding/json"
	"net/http"

	passworddtos "github.com/simon3640/goprojectskeleton/src/application/modules/password/dtos"
	usecases_password "github.com/simon3640/goprojectskeleton/src/application/modules/password/use_cases"
	database "github.com/simon3640/goprojectskeleton/src/infrastructure/databases/goprojectskeleton"
	passwordrepositories "github.com/simon3640/goprojectskeleton/src/infrastructure/databases/goprojectskeleton/repositories/password"
	handlers "github.com/simon3640/goprojectskeleton/src/infrastructure/handlers/shared"
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
func CreatePassword(ctx handlers.HandlerContext) {
	var passwordCreate passworddtos.PasswordCreateNoHash

	if err := json.NewDecoder(*ctx.Body).Decode(&passwordCreate); err != nil {
		http.Error(ctx.ResponseWriter, err.Error(), http.StatusBadRequest)
		return
	}

	passwordRepository := passwordrepositories.NewPasswordRepository(database.GoProjectSkeletondb.DB, providers.Logger)

	ucResult := usecases_password.NewCreatePasswordUseCase(providers.Logger,
		passwordRepository, providers.HashProviderInstance, false,
	).Execute(ctx.Context, ctx.Locale, passwordCreate)
	headers := map[handlers.HTTPHeaderTypeEnum]string{
		handlers.CONTENT_TYPE: string(handlers.APPLICATION_JSON),
	}

	handlers.NewRequestResolver[bool]().ResolveDTO(ctx.ResponseWriter, ucResult, headers)
}
