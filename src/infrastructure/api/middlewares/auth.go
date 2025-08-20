package middlewares

import (
	"gormgoskeleton/src/application/modules/auth"
	app_context "gormgoskeleton/src/application/shared/context"
	"gormgoskeleton/src/application/shared/locales"
	"gormgoskeleton/src/domain/models"
	"gormgoskeleton/src/infrastructure/api"
	database "gormgoskeleton/src/infrastructure/database/gormgoskeleton"
	"gormgoskeleton/src/infrastructure/providers"
	"gormgoskeleton/src/infrastructure/repositories"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("Authorization")

		uc_result := auth.NewAuthUserUseCase(
			providers.Logger,
			repositories.NewUserRepository(database.DB, providers.Logger),
			providers.JWTProviderInstance,
		).Execute(c, locales.EN_US, token)

		if uc_result.HasError() {
			headers := map[api.HTTPHeaderTypeEnum]string{
				api.CONTENT_TYPE: string(api.APPLICATION_JSON),
			}
			content, statusCode := api.NewRequestResolver[models.User]().ResolveDTO(c, uc_result, headers)
			c.JSON(statusCode, content)
			c.Abort()
			return
		}

		user := uc_result.GetData()
		appContext := app_context.NewContextWithUser(user)

		c.Request = c.Request.WithContext(appContext)

		c.Next()

	}
}
