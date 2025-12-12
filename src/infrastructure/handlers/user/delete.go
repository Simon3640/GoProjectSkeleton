package userhandlers

import (
	"net/http"
	"strconv"

	userusecases "github.com/simon3640/goprojectskeleton/src/application/modules/user/use_cases"
	database "github.com/simon3640/goprojectskeleton/src/infrastructure/databases/goprojectskeleton"
	userrepositories "github.com/simon3640/goprojectskeleton/src/infrastructure/databases/goprojectskeleton/repositories/user"
	handlers "github.com/simon3640/goprojectskeleton/src/infrastructure/handlers/shared"
	"github.com/simon3640/goprojectskeleton/src/infrastructure/providers"
)

// DeleteUser delete a user by ID
// @Summary This endpoint Delete a user by ID
// @Description This endpoint Delete a user by ID
// @Tags User
// @Accept json
// @Produce json
// @Param id path int true "ID del usuario"
// @Param Accept-Language header string false "Locale for response messages" Enums(en-US, es-ES) default(en-US)
// @Success 204 {object} nil "Usuario eliminado"
// @Failure 404 {object} map[string]string "Usuario no encontrado"
// @Router /api/user/{id} [delete]
// @Security Bearer
func DeleteUser(ctx handlers.HandlerContext) {
	id, err := strconv.Atoi(ctx.Params["id"])
	if err != nil {
		http.Error(ctx.ResponseWriter, "Invalid ID", http.StatusBadRequest)
		return
	}

	ucResult := userusecases.NewDeleteUserUseCase(providers.Logger,
		userrepositories.NewUserRepository(database.GoProjectSkeletondb.DB, providers.Logger),
	).Execute(ctx.Context, ctx.Locale, uint(id))
	headers := map[handlers.HTTPHeaderTypeEnum]string{
		handlers.CONTENT_TYPE: string(handlers.APPLICATION_JSON),
	}
	handlers.NewRequestResolver[bool]().ResolveDTO(ctx.ResponseWriter, ucResult, headers)
}
