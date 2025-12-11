package userhandlers

import (
	"encoding/json"
	"net/http"

	userdtos "github.com/simon3640/goprojectskeleton/src/application/modules/user/dtos"
	userusecases "github.com/simon3640/goprojectskeleton/src/application/modules/user/use_cases"
	"github.com/simon3640/goprojectskeleton/src/domain/models"
	database "github.com/simon3640/goprojectskeleton/src/infrastructure/databases/goprojectskeleton"
	userrepositories "github.com/simon3640/goprojectskeleton/src/infrastructure/databases/goprojectskeleton/repositories/user"
	handlers "github.com/simon3640/goprojectskeleton/src/infrastructure/handlers/shared"
	"github.com/simon3640/goprojectskeleton/src/infrastructure/providers"
)

// CreateUser
// @Summary This endpoint Create a new user
// @Description This endpoint Create a new user
// @Schemes userdtos.UserCreate
// @Tags User
// @Accept json
// @Produce json
// @Param request body userdtos.UserCreate true "Datos del usuario"
// @Param Accept-Language header string false "Locale for response messages" Enums(en-US, es-ES) default(en-US)
// @Success 201 {object} models.User "Usuario creado"
// @Failure 400 {object} map[string]string "Error de validación"
// @Router /api/user [post]
func CreateUser(ctx handlers.HandlerContext) {
	var userCreate userdtos.UserCreate

	if err := json.NewDecoder(*ctx.Body).Decode(&userCreate); err != nil {
		http.Error(ctx.ResponseWriter, err.Error(), http.StatusBadRequest)
		return
	}

	ucResult := userusecases.NewCreateUserUseCase(providers.Logger,
		userrepositories.NewUserRepository(database.GoProjectSkeletondb.DB, providers.Logger),
	).Execute(ctx.Context, ctx.Locale, userCreate)

	headers := map[handlers.HTTPHeaderTypeEnum]string{
		handlers.CONTENT_TYPE: string(handlers.APPLICATION_JSON),
	}
	handlers.NewRequestResolver[models.User]().ResolveDTO(ctx.ResponseWriter, ucResult, headers)
}
