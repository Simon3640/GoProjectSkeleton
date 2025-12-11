package userhandlers

import (
	"encoding/json"
	"net/http"

	userdtos "github.com/simon3640/goprojectskeleton/src/application/modules/user/dtos"
	userpipes "github.com/simon3640/goprojectskeleton/src/application/modules/user/pipes"
	userusecases "github.com/simon3640/goprojectskeleton/src/application/modules/user/use_cases"
	"github.com/simon3640/goprojectskeleton/src/domain/models"
	database "github.com/simon3640/goprojectskeleton/src/infrastructure/databases/goprojectskeleton"
	userrepositories "github.com/simon3640/goprojectskeleton/src/infrastructure/databases/goprojectskeleton/repositories/user"
	authrepositories "github.com/simon3640/goprojectskeleton/src/infrastructure/databases/goprojectskeleton/repositories/auth"
	handlers "github.com/simon3640/goprojectskeleton/src/infrastructure/handlers/shared"
	"github.com/simon3640/goprojectskeleton/src/infrastructure/providers"
)

// CreateUserAndPassword
// @Summary This endpoint Create a new user
// @Description This endpoint Create a new user and password
// @Schemes models.UserAndPasswordCreate
// @Tags User
// @Accept json
// @Produce json
// @Param request body dtos.UserAndPasswordCreate true "Datos del usuario"
// @Param Accept-Language header string false "Locale for response messages" Enums(en-US, es-ES) default(en-US)
// @Success 201 {object} models.User "Usuario creado"
// @Failure 400 {object} map[string]string "Error de validación"
// @Router /api/user-password [post]
func CreateUserAndPassword(ctx handlers.HandlerContext) {
	var userCreate userdtos.UserAndPasswordCreate

	if err := json.NewDecoder(*ctx.Body).Decode(&userCreate); err != nil {
		http.Error(ctx.ResponseWriter, err.Error(), http.StatusBadRequest)
		return
	}
	createUserSendEmailUC := userusecases.NewCreateUserSendEmailUseCase(
		providers.Logger,
		providers.HashProviderInstance,
		authrepositories.NewOneTimeTokenRepository(database.GoProjectSkeletondb.DB, providers.Logger),
	)

	createUserPasswordUC := userusecases.NewCreateUserAndPasswordUseCase(providers.Logger,
		userrepositories.NewUserRepository(database.GoProjectSkeletondb.DB, providers.Logger),
		providers.HashProviderInstance,
	)
	ucResult := userpipes.NewCreateUserPipe(ctx.Context,
		ctx.Locale,
		createUserPasswordUC,
		createUserSendEmailUC,
	).Execute(userCreate)

	headers := map[handlers.HTTPHeaderTypeEnum]string{
		handlers.CONTENT_TYPE: string(handlers.APPLICATION_JSON),
	}
	handlers.NewRequestResolver[models.User]().ResolveDTO(ctx.ResponseWriter, ucResult, headers)
}
