package routes

import (
	"go/types"
	"net/http"
	"strconv"

	user_pipes "gormgoskeleton/src/application/modules/user/pipes"
	usecases_user "gormgoskeleton/src/application/modules/user/use_cases"
	"gormgoskeleton/src/application/shared/locales"
	"gormgoskeleton/src/domain/models"
	domain_utils "gormgoskeleton/src/domain/utils"
	"gormgoskeleton/src/infrastructure/api"
	database "gormgoskeleton/src/infrastructure/database/gormgoskeleton"
	"gormgoskeleton/src/infrastructure/providers"
	"gormgoskeleton/src/infrastructure/repositories"

	"github.com/gin-gonic/gin"
)

// CreateUser
// @Summary This endpoint Create a new user
// @Description This endpoint Create a new user
// @Schemes models.UserCreate
// @Tags User
// @Accept json
// @Produce json
// @Param request body models.UserCreate true "Datos del usuario"
// @Success 201 {object} models.User "Usuario creado"
// @Failure 400 {object} map[string]string "Error de validación"
// @Router /api/user [post]
func createUser(c *gin.Context) {
	var userCreate models.UserCreate

	if err := c.ShouldBindJSON(&userCreate); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	uc_result := usecases_user.NewCreateUserUseCase(providers.Logger,
		repositories.NewUserRepository(database.DB, providers.Logger),
	).Execute(c.Request.Context(), locales.EN_US, userCreate)
	headers := map[api.HTTPHeaderTypeEnum]string{
		api.CONTENT_TYPE: string(api.APPLICATION_JSON),
	}
	content, statusCode := api.NewRequestResolver[models.User]().ResolveDTO(c, uc_result, headers)

	c.JSON(statusCode, content)
}

// GetUser
// @Summary This endpoint Get a user by ID
// @Description This endpoint Get a user by ID
// @Tags User
// @Accept json
// @Produce json
// @Param id path int true "ID del usuario"
// @Success 200 {object} models.User "Usuario"
// @Failure 404 {object} map[string]string "Usuario no encontrado"
// @Router /api/user/{id} [get]
// @Security Bearer
func getUser(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	uc_result := usecases_user.NewGetUserUseCase(providers.Logger,
		repositories.NewUserRepository(database.DB, providers.Logger),
	).Execute(c.Request.Context(), locales.EN_US, uint(id))
	headers := map[api.HTTPHeaderTypeEnum]string{
		api.CONTENT_TYPE: string(api.APPLICATION_JSON),
	}
	content, statusCode := api.NewRequestResolver[models.User]().ResolveDTO(c, uc_result, headers)

	c.JSON(statusCode, content)
}

// UpdateUser
// @Summary This endpoint Update a user by ID
// @Description This endpoint Update a user by ID
// @Tags User
// @Accept json
// @Produce json
// @Param id path int true "ID del usuario"
// @Param request body models.UserUpdateBase true "Datos del usuario"
// @Success 200 {object} models.User "Usuario actualizado"
// @Failure 400 {object} map[string]string "Error de validación"
// @Router /api/user/{id} [patch]
// @Security Bearer
func updateUser(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var userUpdate models.UserUpdate

	if err := c.ShouldBindJSON(&userUpdate); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userUpdate.ID = uint(id)
	uc_result := usecases_user.NewUpdateUserUseCase(providers.Logger,
		repositories.NewUserRepository(database.DB, providers.Logger),
	).Execute(c.Request.Context(), locales.EN_US, userUpdate)
	headers := map[api.HTTPHeaderTypeEnum]string{
		api.CONTENT_TYPE: string(api.APPLICATION_JSON),
	}
	content, statusCode := api.NewRequestResolver[models.User]().ResolveDTO(c, uc_result, headers)

	c.JSON(statusCode, content)
}

// DeleteUser
// @Summary This endpoint Delete a user by ID
// @Description This endpoint Delete a user by ID
// @Tags User
// @Accept json
// @Produce json
// @Param id path int true "ID del usuario"
// @Success 204 {object} nil "Usuario eliminado"
// @Failure 404 {object} map[string]string "Usuario no encontrado"
// @Router /api/user/{id} [delete]
// @Security Bearer
func deleteUser(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	uc_result := usecases_user.NewDeleteUserUseCase(providers.Logger,
		repositories.NewUserRepository(database.DB, providers.Logger),
	).Execute(c.Request.Context(), locales.EN_US, uint(id))
	headers := map[api.HTTPHeaderTypeEnum]string{
		api.CONTENT_TYPE: string(api.APPLICATION_JSON),
	}
	content, statusCode := api.NewRequestResolver[types.Nil]().ResolveDTO(c, uc_result, headers)

	c.JSON(statusCode, content)
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
//
// @Success 200 {array} models.User "List of users"
// @Failure 400 {object} map[string]string "Bad request"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /api/user [get]
func getAllUser(c *gin.Context) {
	queryParams, exists := c.Get("queryParams")
	if !exists {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Query parameters not found"})
		return
	}
	uc_result := usecases_user.NewGetAllUserUseCase(providers.Logger,
		repositories.NewUserRepository(database.DB, providers.Logger),
	).Execute(c.Request.Context(), locales.EN_US, queryParams.(domain_utils.QueryPayloadBuilder[models.User]))
	headers := map[api.HTTPHeaderTypeEnum]string{
		api.CONTENT_TYPE: string(api.APPLICATION_JSON),
	}
	content, statusCode := api.NewRequestResolver[[]models.User]().ResolveDTO(c, uc_result, headers)

	c.JSON(statusCode, content)
}

// CreateUserAndPassword
// @Summary This endpoint Create a new user
// @Description This endpoint Create a new user and password
// @Schemes models.UserAndPasswordCreate
// @Tags User
// @Accept json
// @Produce json
// @Param request body models.UserAndPasswordCreate true "Datos del usuario"
// @Success 201 {object} models.User "Usuario creado"
// @Failure 400 {object} map[string]string "Error de validación"
// @Router /api/user-password [post]
func createUserAndPassword(c *gin.Context) {
	var userCreate models.UserAndPasswordCreate

	if err := c.ShouldBindJSON(&userCreate); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	userRepository := repositories.NewUserRepository(database.DB, providers.Logger)
	hashProvider := providers.HashProviderInstance

	uc_create_user_email := usecases_user.NewCreateUserSendEmailUseCase(
		providers.Logger,
		providers.JWTProviderInstance,
	)

	uc_create_user_password := usecases_user.NewCreateUserAndPasswordUseCase(providers.Logger,
		userRepository,
		hashProvider,
	)

	uc_result := user_pipes.NewCreateUserPipe(c.Request.Context(),
		locales.EN_US,
		uc_create_user_password,
		uc_create_user_email,
	).Execute(userCreate)

	// uc_result := usecases_user.NewCreateUserAndPasswordUseCase(providers.Logger,
	// 	userRepository,
	// 	hashProvider,
	// ).Execute(c.Request.Context(), locales.EN_US, userCreate)
	headers := map[api.HTTPHeaderTypeEnum]string{
		api.CONTENT_TYPE: string(api.APPLICATION_JSON),
	}
	content, statusCode := api.NewRequestResolver[models.User]().ResolveDTO(c, uc_result, headers)

	c.JSON(statusCode, content)
}
