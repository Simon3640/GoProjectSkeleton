package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	userpipes "github.com/simon3640/goprojectskeleton/src/application/modules/user/pipes"
	userusecases "github.com/simon3640/goprojectskeleton/src/application/modules/user/use_cases"
	dtos "github.com/simon3640/goprojectskeleton/src/application/shared/DTOs"
	"github.com/simon3640/goprojectskeleton/src/domain/models"
	domain_utils "github.com/simon3640/goprojectskeleton/src/domain/utils"
	database "github.com/simon3640/goprojectskeleton/src/infrastructure/database/goprojectskeleton"
	"github.com/simon3640/goprojectskeleton/src/infrastructure/providers"
	"github.com/simon3640/goprojectskeleton/src/infrastructure/repositories"
)

// CreateUser
// @Summary This endpoint Create a new user
// @Description This endpoint Create a new user
// @Schemes dtos.UserCreate
// @Tags User
// @Accept json
// @Produce json
// @Param request body dtos.UserCreate true "Datos del usuario"
// @Param Accept-Language header string false "Locale for response messages" Enums(en-US, es-ES) default(en-US)
// @Success 201 {object} models.User "Usuario creado"
// @Failure 400 {object} map[string]string "Error de validación"
// @Router /api/user [post]
func CreateUser(ctx HandlerContext) {
	var userCreate dtos.UserCreate

	if err := json.NewDecoder(*ctx.Body).Decode(&userCreate); err != nil {
		http.Error(ctx.ResponseWriter, err.Error(), http.StatusBadRequest)
		return
	}

	ucResult := userusecases.NewCreateUserUseCase(providers.Logger,
		repositories.NewUserRepository(database.GoProjectSkeletondb.DB, providers.Logger),
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
// @Param Accept-Language header string false "Locale for response messages" Enums(en-US, es-ES) default(en-US)
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

	ucResult := userusecases.NewGetUserUseCase(providers.Logger,
		repositories.NewUserRepository(database.GoProjectSkeletondb.DB, providers.Logger),
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
// @Param Accept-Language header string false "Locale for response messages" Enums(en-US, es-ES) default(en-US)
// @Param request body models.UserUpdateBase true "Datos del usuario"
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
	ucResult := userusecases.NewUpdateUserUseCase(providers.Logger,
		repositories.NewUserRepository(database.GoProjectSkeletondb.DB, providers.Logger),
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
// @Param Accept-Language header string false "Locale for response messages" Enums(en-US, es-ES) default(en-US)
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

	ucResult := userusecases.NewDeleteUserUseCase(providers.Logger,
		repositories.NewUserRepository(database.GoProjectSkeletondb.DB, providers.Logger),
	).Execute(ctx.c, ctx.Locale, uint(id))
	headers := map[HTTPHeaderTypeEnum]string{
		CONTENT_TYPE: string(APPLICATION_JSON),
	}
	NewRequestResolver[bool]().ResolveDTO(ctx.ResponseWriter, ucResult, headers)
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
// @Param Accept-Language header string false "Locale for response messages" Enums(en-US, es-ES) default(en-US)
//
// @Success 200 {object} dtos.UserMultiResponse "List of users"
// @Failure 400 {object} map[string]string "Bad request"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /api/user [get]
func GetAllUser(ctx HandlerContext) {
	queryParams := domain_utils.NewQueryPayloadBuilder[models.User](ctx.Query.Sorts, ctx.Query.Filters, ctx.Query.Page, ctx.Query.PageSize)
	ucResult := userusecases.NewGetAllUserUseCase(providers.Logger,
		repositories.NewUserRepository(database.GoProjectSkeletondb.DB, providers.Logger),
		providers.CacheProviderInstance,
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
// @Param Accept-Language header string false "Locale for response messages" Enums(en-US, es-ES) default(en-US)
// @Success 201 {object} models.User "Usuario creado"
// @Failure 400 {object} map[string]string "Error de validación"
// @Router /api/user-password [post]
func CreateUserAndPassword(ctx HandlerContext) {
	var userCreate dtos.UserAndPasswordCreate

	if err := json.NewDecoder(*ctx.Body).Decode(&userCreate); err != nil {
		http.Error(ctx.ResponseWriter, err.Error(), http.StatusBadRequest)
		return
	}
	createUserSendEmailUC := userusecases.NewCreateUserSendEmailUseCase(
		providers.Logger,
		providers.HashProviderInstance,
		repositories.NewOneTimeTokenRepository(database.GoProjectSkeletondb.DB, providers.Logger),
	)

	createUserPasswordUC := userusecases.NewCreateUserAndPasswordUseCase(providers.Logger,
		repositories.NewUserRepository(database.GoProjectSkeletondb.DB, providers.Logger),
		providers.HashProviderInstance,
	)
	ucResult := userpipes.NewCreateUserPipe(ctx.c,
		ctx.Locale,
		createUserPasswordUC,
		createUserSendEmailUC,
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
// @Param Accept-Language header string false "Locale for response messages" Enums(en-US, es-ES) default(en-US)
// @Success 200 {object} bool "Usuario activado"
// @Failure 400 {object} map[string]string "Error de validación"
// @Router /api/user/activate [post]
func ActivateUser(ctx HandlerContext) {
	var userActivate dtos.UserActivate

	if err := json.NewDecoder(*ctx.Body).Decode(&userActivate); err != nil {
		http.Error(ctx.ResponseWriter, err.Error(), http.StatusBadRequest)
		return
	}

	ucResult := userusecases.NewActivateUserUseCase(
		providers.Logger,
		repositories.NewUserRepository(database.GoProjectSkeletondb.DB, providers.Logger),
		repositories.NewOneTimeTokenRepository(database.GoProjectSkeletondb.DB, providers.Logger),
		providers.HashProviderInstance,
	).Execute(ctx.c, ctx.Locale, userActivate)
	headers := map[HTTPHeaderTypeEnum]string{
		CONTENT_TYPE: string(APPLICATION_JSON),
	}
	NewRequestResolver[bool]().ResolveDTO(ctx.ResponseWriter, ucResult, headers)
}

// ResendWelcomeEmail Resend welcome email to user
// @Summary Resend welcome email to user
// @Description This endpoint resends the welcome email with a new activation token to the user
// @Tags User
// @Accept json
// @Produce json
// @Param request body dtos.ResendWelcomeEmailRequest true "Email del usuario"
// @Param Accept-Language header string false "Locale for response messages" Enums(en-US, es-ES) default(en-US)
// @Success 200 {object} bool "Correo de bienvenida reenviado"
// @Failure 400 {object} map[string]string "Error de validación"
// @Failure 404 {object} map[string]string "Usuario no encontrado"
// @Router /api/user/resend-welcome-email [post]
func ResendWelcomeEmail(ctx HandlerContext) {
	var resendRequest dtos.ResendWelcomeEmailRequest

	if err := json.NewDecoder(*ctx.Body).Decode(&resendRequest); err != nil {
		http.Error(ctx.ResponseWriter, err.Error(), http.StatusBadRequest)
		return
	}

	ucResult := userusecases.NewResendWelcomeEmailUseCase(
		providers.Logger,
		providers.HashProviderInstance,
		repositories.NewUserRepository(database.GoProjectSkeletondb.DB, providers.Logger),
		repositories.NewOneTimeTokenRepository(database.GoProjectSkeletondb.DB, providers.Logger),
	).Execute(ctx.c, ctx.Locale, resendRequest)
	headers := map[HTTPHeaderTypeEnum]string{
		CONTENT_TYPE: string(APPLICATION_JSON),
	}
	NewRequestResolver[bool]().ResolveDTO(ctx.ResponseWriter, ucResult, headers)
}
