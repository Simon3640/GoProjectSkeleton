package handlers

import (
	"encoding/json"
	"go/types"
	"net/http"
	"strconv"

	user_pipes "gormgoskeleton/src/application/modules/user/pipes"
	usecases_user "gormgoskeleton/src/application/modules/user/use_cases"
	dtos "gormgoskeleton/src/application/shared/DTOs"
	"gormgoskeleton/src/domain/models"
	domain_utils "gormgoskeleton/src/domain/utils"
	database "gormgoskeleton/src/infrastructure/database/gormgoskeleton"
	"gormgoskeleton/src/infrastructure/providers"
	"gormgoskeleton/src/infrastructure/repositories"
)

// CreateUser
// @Summary This endpoint Create a new user
// @Description This endpoint Create a new user
// @Schemes dtos.UserCreate
// @Tags User
// @Accept json
// @Produce json
// @Param request body dtos.UserCreate true "Datos del usuario"
// @Param X-Locale header string false "Locale for response messages" Enums(en-US, es-ES) default(en-US)
// @Success 201 {object} models.User "Usuario creado"
// @Failure 400 {object} map[string]string "Error de validación"
// @Router /api/user [post]
func CreateUser(ctx HandlerContext) {
	var userCreate dtos.UserCreate

	if err := json.NewDecoder(*ctx.Body).Decode(&userCreate); err != nil {
		http.Error(ctx.ResponseWriter, err.Error(), http.StatusBadRequest)
		return
	}

	ucResult := usecases_user.NewCreateUserUseCase(providers.Logger,
		repositories.NewUserRepository(database.DB, providers.Logger),
	).Execute(ctx.c, ctx.Locale, userCreate)

	headers := map[HTTPHeaderTypeEnum]string{
		CONTENT_TYPE: string(APPLICATION_JSON),
	}
	NewRequestResolver[models.User]().ResolveDTO(ctx.ResponseWriter, ucResult, headers)

}

// GetUser
// @Summary This endpoint Get a user by ID
// @Description This endpoint Get a user by ID
// @Tags User
// @Accept json
// @Produce json
// @Param id path int true "ID del usuario"
// @Param X-Locale header string false "Locale for response messages" Enums(en-US, es-ES) default(en-US)
// @Success 200 {object} models.User "Usuario"
// @Failure 404 {object} map[string]string "Usuario no encontrado"
// @Router /api/user/{id} [get]
// @Security Bearer
func GetUser(ctx HandlerContext) {
	id, err := strconv.Atoi(ctx.Params["id"])
	if err != nil {
		http.Error(ctx.ResponseWriter, "Invalid ID", http.StatusBadRequest)
		return
	}

	ucResult := usecases_user.NewGetUserUseCase(providers.Logger,
		repositories.NewUserRepository(database.DB, providers.Logger),
	).Execute(ctx.c, ctx.Locale, uint(id))
	headers := map[HTTPHeaderTypeEnum]string{
		CONTENT_TYPE: string(APPLICATION_JSON),
	}
	NewRequestResolver[models.User]().ResolveDTO(ctx.ResponseWriter, ucResult, headers)
}

// UpdateUser
// @Summary This endpoint Update a user by ID
// @Description This endpoint Update a user by ID
// @Tags User
// @Accept json
// @Produce json
// @Param id path int true "ID del usuario"
// @Param X-Locale header string false "Locale for response messages" Enums(en-US, es-ES) default(en-US)
// @Param request body dtos.UserUpdateBase true "Datos del usuario"
// @Success 200 {object} models.User "Usuario actualizado"
// @Failure 400 {object} map[string]string "Error de validación"
// @Router /api/user/{id} [patch]
// @Security Bearer
func UpdateUser(ctx HandlerContext) {
	id, err := strconv.Atoi(ctx.Params["id"])
	if err != nil {
		http.Error(ctx.ResponseWriter, "Invalid ID", http.StatusBadRequest)
		return
	}

	var userUpdate dtos.UserUpdate
	if err := json.NewDecoder(*ctx.Body).Decode(&userUpdate); err != nil {
		http.Error(ctx.ResponseWriter, err.Error(), http.StatusBadRequest)
		return
	}

	userUpdate.ID = uint(id)
	ucResult := usecases_user.NewUpdateUserUseCase(providers.Logger,
		repositories.NewUserRepository(database.DB, providers.Logger),
	).Execute(ctx.c, ctx.Locale, userUpdate)
	headers := map[HTTPHeaderTypeEnum]string{
		CONTENT_TYPE: string(APPLICATION_JSON),
	}
	NewRequestResolver[models.User]().ResolveDTO(ctx.ResponseWriter, ucResult, headers)
}

// DeleteUser
// @Summary This endpoint Delete a user by ID
// @Description This endpoint Delete a user by ID
// @Tags User
// @Accept json
// @Produce json
// @Param id path int true "ID del usuario"
// @Param X-Locale header string false "Locale for response messages" Enums(en-US, es-ES) default(en-US)
// @Success 204 {object} nil "Usuario eliminado"
// @Failure 404 {object} map[string]string "Usuario no encontrado"
// @Router /api/user/{id} [delete]
// @Security Bearer
func DeleteUser(ctx HandlerContext) {
	id, err := strconv.Atoi(ctx.Params["id"])
	if err != nil {
		http.Error(ctx.ResponseWriter, "Invalid ID", http.StatusBadRequest)
		return
	}

	ucResult := usecases_user.NewDeleteUserUseCase(providers.Logger,
		repositories.NewUserRepository(database.DB, providers.Logger),
	).Execute(ctx.c, ctx.Locale, uint(id))
	headers := map[HTTPHeaderTypeEnum]string{
		CONTENT_TYPE: string(APPLICATION_JSON),
	}
	NewRequestResolver[types.Nil]().ResolveDTO(ctx.ResponseWriter, ucResult, headers)
}

// GetAllUser
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
// @Param X-Locale header string false "Locale for response messages" Enums(en-US, es-ES) default(en-US)
//
// @Success 200 {object} dtos.UserMultiResponse "List of users"
// @Failure 400 {object} map[string]string "Bad request"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /api/user [get]
func GetAllUser(ctx HandlerContext) {
	queryParams := domain_utils.NewQueryPayloadBuilder[models.User](ctx.Query.Sorts, ctx.Query.Filters, ctx.Query.Page, ctx.Query.PageSize)
	ucResult := usecases_user.NewGetAllUserUseCase(providers.Logger,
		repositories.NewUserRepository(database.DB, providers.Logger),
	).Execute(ctx.c, ctx.Locale, queryParams)
	headers := map[HTTPHeaderTypeEnum]string{
		CONTENT_TYPE: string(APPLICATION_JSON),
	}
	NewRequestResolver[dtos.UserMultiResponse]().ResolveDTO(ctx.ResponseWriter, ucResult, headers)
}

// CreateUserAndPassword
// @Summary This endpoint Create a new user
// @Description This endpoint Create a new user and password
// @Schemes models.UserAndPasswordCreate
// @Tags User
// @Accept json
// @Produce json
// @Param request body dtos.UserAndPasswordCreate true "Datos del usuario"
// @Param X-Locale header string false "Locale for response messages" Enums(en-US, es-ES) default(en-US)
// @Success 201 {object} models.User "Usuario creado"
// @Failure 400 {object} map[string]string "Error de validación"
// @Router /api/user-password [post]
func CreateUserAndPassword(ctx HandlerContext) {
	var userCreate dtos.UserAndPasswordCreate

	if err := json.NewDecoder(*ctx.Body).Decode(&userCreate); err != nil {
		http.Error(ctx.ResponseWriter, err.Error(), http.StatusBadRequest)
		return
	}
	uc_create_user_email := usecases_user.NewCreateUserSendEmailUseCase(
		providers.Logger,
		providers.HashProviderInstance,
		repositories.NewOneTimeTokenRepository(database.DB, providers.Logger),
	)

	uc_create_user_password := usecases_user.NewCreateUserAndPasswordUseCase(providers.Logger,
		repositories.NewUserRepository(database.DB, providers.Logger),
		providers.HashProviderInstance,
	)
	ucResult := user_pipes.NewCreateUserPipe(ctx.c,
		ctx.Locale,
		uc_create_user_password,
		uc_create_user_email,
	).Execute(userCreate)

	headers := map[HTTPHeaderTypeEnum]string{
		CONTENT_TYPE: string(APPLICATION_JSON),
	}
	NewRequestResolver[models.User]().ResolveDTO(ctx.ResponseWriter, ucResult, headers)
}

// ActivateUser
// @Summary This endpoint Activate a user by token
// @Description This endpoint Activate a user by token
// @Schemes models.UserActivate
// @Tags User
// @Accept json
// @Produce json
// @Param request body dtos.UserActivate true "Token de activación"
// @Param X-Locale header string false "Locale for response messages" Enums(en-US, es-ES) default(en-US)
// @Success 200 {object} bool "Usuario activado"
// @Failure 400 {object} map[string]string "Error de validación"
// @Router /api/user/activate [post]
func ActivateUser(ctx HandlerContext) {
	var userActivate dtos.UserActivate

	if err := json.NewDecoder(*ctx.Body).Decode(&userActivate); err != nil {
		http.Error(ctx.ResponseWriter, err.Error(), http.StatusBadRequest)
		return
	}

	ucResult := usecases_user.NewActivateUserUseCase(
		providers.Logger,
		repositories.NewUserRepository(database.DB, providers.Logger),
		repositories.NewOneTimeTokenRepository(database.DB, providers.Logger),
		providers.HashProviderInstance,
	).Execute(ctx.c, ctx.Locale, userActivate)
	headers := map[HTTPHeaderTypeEnum]string{
		CONTENT_TYPE: string(APPLICATION_JSON),
	}
	NewRequestResolver[bool]().ResolveDTO(ctx.ResponseWriter, ucResult, headers)
}
