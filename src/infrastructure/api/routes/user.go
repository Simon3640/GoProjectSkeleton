package routes

import (
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
