package routes

import (
	"go/types"
	"net/http"
	"strconv"

	"gormgoskeleton/src/application/shared/locales"
	usecases_user "gormgoskeleton/src/application/use_cases/user"
	"gormgoskeleton/src/domain/models"
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
		repositories.NewUserRepository(database.DB),
	).Execute(c, locales.EN_US, userCreate)
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
func getUser(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	uc_result := usecases_user.NewGetUserUseCase(providers.Logger,
		repositories.NewUserRepository(database.DB),
	).Execute(c, locales.EN_US, id)
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

	userUpdate.ID = id

	uc_result := usecases_user.NewUpdateUserUseCase(providers.Logger,
		repositories.NewUserRepository(database.DB),
	).Execute(c, locales.EN_US, userUpdate)
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
func deleteUser(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	uc_result := usecases_user.NewDeleteUserUseCase(providers.Logger,
		repositories.NewUserRepository(database.DB),
	).Execute(c, locales.EN_US, id)
	headers := map[api.HTTPHeaderTypeEnum]string{
		api.CONTENT_TYPE: string(api.APPLICATION_JSON),
	}
	content, statusCode := api.NewRequestResolver[types.Nil]().ResolveDTO(c, uc_result, headers)

	c.JSON(statusCode, content)
}

// GetAllUser
// @Summary This endpoint Get all users
// @Description This endpoint Get all users
// @Tags User
// @Accept json
// @Produce json
// @Success 200 {object} []models.User "Usuarios"
// @Router /api/user [get]
func getAllUser(c *gin.Context) {
	uc_result := usecases_user.NewGetAllUserUseCase(providers.Logger,
		repositories.NewUserRepository(database.DB),
	).Execute(c, locales.EN_US, usecases_user.Nil{})
	headers := map[api.HTTPHeaderTypeEnum]string{
		api.CONTENT_TYPE: string(api.APPLICATION_JSON),
	}
	content, statusCode := api.NewRequestResolver[[]models.User]().ResolveDTO(c, uc_result, headers)

	c.JSON(statusCode, content)
}
