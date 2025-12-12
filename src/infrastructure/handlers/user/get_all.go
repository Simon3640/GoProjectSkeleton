package userhandlers

import (
	userdtos "github.com/simon3640/goprojectskeleton/src/application/modules/user/dtos"
	userusecases "github.com/simon3640/goprojectskeleton/src/application/modules/user/use_cases"
	"github.com/simon3640/goprojectskeleton/src/domain/models"
	domainutils "github.com/simon3640/goprojectskeleton/src/domain/utils"
	database "github.com/simon3640/goprojectskeleton/src/infrastructure/databases/goprojectskeleton"
	userrepositories "github.com/simon3640/goprojectskeleton/src/infrastructure/databases/goprojectskeleton/repositories/user"
	handlers "github.com/simon3640/goprojectskeleton/src/infrastructure/handlers/shared"
	"github.com/simon3640/goprojectskeleton/src/infrastructure/providers"
)

// GetAllUser get all users
// @Summary Get all users
// @Description Retrieve all users with support for filtering, sorting, and pagination.
// @Tags User
// @Accept json
// @Produce json
// @Security Bearer
//
// @Param filter query []string false "Filter users in the format column:operator:value (e.g. Name:eq:Admin, Age:gt:18)"
// @Param sort query []string false "Sort users in the format column:asc|desc (e.g. Name:asc, CreatedAt:desc)"
// @Param page query int false "Page number (default: 1)"
// @Param page_size query int false "Number of items per page (default: 10)"
// @Param Accept-Language header string false "Locale for response messages" Enums(en-US, es-ES) default(en-US)
//
// @Success 200 {object} userdtos.UserMultiResponse "List of users"
// @Failure 400 {object} map[string]string "Bad request"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /api/user [get]
func GetAllUser(ctx handlers.HandlerContext) {
	queryParams := domainutils.NewQueryPayloadBuilder[models.User](ctx.Query.Sorts, ctx.Query.Filters, ctx.Query.Page, ctx.Query.PageSize)
	ucResult := userusecases.NewGetAllUserUseCase(providers.Logger,
		userrepositories.NewUserRepository(database.GoProjectSkeletondb.DB, providers.Logger),
		providers.CacheProviderInstance,
	).Execute(ctx.Context, ctx.Locale, queryParams)
	headers := map[handlers.HTTPHeaderTypeEnum]string{
		handlers.CONTENT_TYPE: string(handlers.APPLICATION_JSON),
	}
	handlers.NewRequestResolver[userdtos.UserMultiResponse]().ResolveDTO(ctx.ResponseWriter, ucResult, headers)
}
