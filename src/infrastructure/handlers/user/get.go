package userhandlers

import (
	"net/http"
	"strconv"

	userusecases "github.com/simon3640/goprojectskeleton/src/application/modules/user/use_cases"
	"github.com/simon3640/goprojectskeleton/src/application/shared/observability"
	usecase "github.com/simon3640/goprojectskeleton/src/application/shared/use_case"
	usermodels "github.com/simon3640/goprojectskeleton/src/domain/user/models"
	database "github.com/simon3640/goprojectskeleton/src/infrastructure/databases/goprojectskeleton"
	userrepositories "github.com/simon3640/goprojectskeleton/src/infrastructure/databases/goprojectskeleton/repositories/user"
	handlers "github.com/simon3640/goprojectskeleton/src/infrastructure/handlers/shared"
	"github.com/simon3640/goprojectskeleton/src/infrastructure/providers"
)

// GetUser get a user by ID
// @Summary This endpoint Get a user by ID
// @Description This endpoint Get a user by ID
// @Tags User
// @Accept json
// @Produce json
// @Param id path int true "ID del usuario"
// @Param Accept-Language header string false "Locale for response messages" Enums(en-US, es-ES) default(en-US)
// @Success 200 {object} usermodels.User "Usuario"
// @Failure 404 {object} map[string]string "Usuario no encontrado"
// @Router /api/user/{id} [get]
// @Security Bearer
func GetUser(ctx handlers.HandlerContext) {
	id, err := strconv.Atoi(ctx.Params["id"])
	if err != nil {
		http.Error(ctx.ResponseWriter, "Invalid ID", http.StatusBadRequest)
		return
	}

	uc := userusecases.NewGetUserUseCase(
		userrepositories.NewUserRepository(database.GoProjectSkeletondb.DB, providers.Logger),
	)
	ucResult := usecase.InstrumentUseCase(
		uc,
		ctx.Context,
		ctx.Locale,
		uint(id),
		observability.GetObservabilityComponents().Tracer,
		observability.GetObservabilityComponents().Metrics,
		observability.GetObservabilityComponents().Clock,
		"get_user_use_case",
	)
	headers := map[handlers.HTTPHeaderTypeEnum]string{
		handlers.CONTENT_TYPE: string(handlers.APPLICATION_JSON),
	}
	handlers.NewRequestResolver[usermodels.User]().ResolveDTO(ctx.ResponseWriter, ucResult, headers)
}
