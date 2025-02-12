package routes

import (
	"net/http"

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
