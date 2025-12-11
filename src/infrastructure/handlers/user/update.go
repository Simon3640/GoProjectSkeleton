package userhandlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	userdtos "github.com/simon3640/goprojectskeleton/src/application/modules/user/dtos"
	userusecases "github.com/simon3640/goprojectskeleton/src/application/modules/user/use_cases"
	"github.com/simon3640/goprojectskeleton/src/domain/models"
	database "github.com/simon3640/goprojectskeleton/src/infrastructure/databases/goprojectskeleton"
	userrepositories "github.com/simon3640/goprojectskeleton/src/infrastructure/databases/goprojectskeleton/repositories/user"
	handlers "github.com/simon3640/goprojectskeleton/src/infrastructure/handlers/shared"
	"github.com/simon3640/goprojectskeleton/src/infrastructure/providers"
)

// UpdateUser
// @Summary This endpoint Update a user by ID
// @Description This endpoint Update a user by ID
// @Tags User
// @Accept json
// @Produce json
// @Param id path int true "ID del usuario"
// @Param Accept-Language header string false "Locale for response messages" Enums(en-US, es-ES) default(en-US)
// @Param request body userdtos.UserUpdate true "Datos del usuario"
// @Success 200 {object} models.User "Usuario actualizado"
// @Failure 400 {object} map[string]string "Error de validación"
// @Router /api/user/{id} [patch]
// @Security Bearer
func UpdateUser(ctx handlers.HandlerContext) {
	id, err := strconv.Atoi(ctx.Params["id"])
	if err != nil {
		http.Error(ctx.ResponseWriter, "Invalid ID", http.StatusBadRequest)
		return
	}

	var userUpdate userdtos.UserUpdate
	if err := json.NewDecoder(*ctx.Body).Decode(&userUpdate); err != nil {
		http.Error(ctx.ResponseWriter, err.Error(), http.StatusBadRequest)
		return
	}

	userUpdate.ID = uint(id)
	ucResult := userusecases.NewUpdateUserUseCase(providers.Logger,
		userrepositories.NewUserRepository(database.GoProjectSkeletondb.DB, providers.Logger),
	).Execute(ctx.Context, ctx.Locale, userUpdate)
	headers := map[handlers.HTTPHeaderTypeEnum]string{
		handlers.CONTENT_TYPE: string(handlers.APPLICATION_JSON),
	}
	handlers.NewRequestResolver[models.User]().ResolveDTO(ctx.ResponseWriter, ucResult, headers)
}
