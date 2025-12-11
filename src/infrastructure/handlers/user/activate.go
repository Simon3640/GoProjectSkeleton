// Package userhandlers contains the handlers for the user module
package userhandlers

import (
	"encoding/json"
	"net/http"

	userdtos "github.com/simon3640/goprojectskeleton/src/application/modules/user/dtos"
	userusecases "github.com/simon3640/goprojectskeleton/src/application/modules/user/use_cases"
	database "github.com/simon3640/goprojectskeleton/src/infrastructure/databases/goprojectskeleton"
	authrepositories "github.com/simon3640/goprojectskeleton/src/infrastructure/databases/goprojectskeleton/repositories/auth"
	userrepositories "github.com/simon3640/goprojectskeleton/src/infrastructure/databases/goprojectskeleton/repositories/user"
	handlers "github.com/simon3640/goprojectskeleton/src/infrastructure/handlers/shared"
	"github.com/simon3640/goprojectskeleton/src/infrastructure/providers"
)

// ActivateUser activate a user by token
// @Summary This endpoint Activate a user by token
// @Description This endpoint Activate a user by token
// @Schemes userdtos.UserActivate
// @Tags User
// @Accept json
// @Produce json
// @Param request body userdtos.UserActivate true "Token de activación"
// @Param Accept-Language header string false "Locale for response messages" Enums(en-US, es-ES) default(en-US)
// @Success 200 {object} bool "Usuario activado"
// @Failure 400 {object} map[string]string "Error de validación"
// @Router /api/user/activate [post]
func ActivateUser(ctx handlers.HandlerContext) {
	var userActivate userdtos.UserActivate

	if err := json.NewDecoder(*ctx.Body).Decode(&userActivate); err != nil {
		http.Error(ctx.ResponseWriter, err.Error(), http.StatusBadRequest)
		return
	}

	ucResult := userusecases.NewActivateUserUseCase(
		providers.Logger,
		userrepositories.NewUserRepository(database.GoProjectSkeletondb.DB, providers.Logger),
		authrepositories.NewOneTimeTokenRepository(database.GoProjectSkeletondb.DB, providers.Logger),
		providers.HashProviderInstance,
	).Execute(ctx.Context, ctx.Locale, userActivate)
	headers := map[handlers.HTTPHeaderTypeEnum]string{
		handlers.CONTENT_TYPE: string(handlers.APPLICATION_JSON),
	}
	handlers.NewRequestResolver[bool]().ResolveDTO(ctx.ResponseWriter, ucResult, headers)
}
