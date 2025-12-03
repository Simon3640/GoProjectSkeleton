package middlewares

import (
	"context"
	api "goprojectskeleton/gin"
	"goprojectskeleton/src/application/modules/auth"
	app_context "goprojectskeleton/src/application/shared/context"
	"goprojectskeleton/src/application/shared/locales"
	"goprojectskeleton/src/domain/models"
	database "goprojectskeleton/src/infrastructure/database/goprojectskeleton"
	"goprojectskeleton/src/infrastructure/providers"
	"goprojectskeleton/src/infrastructure/repositories"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("Authorization")

		uc_result := auth.NewAuthUserUseCase(
			providers.Logger,
			repositories.NewUserRepository(database.GoProjectSkeletondb.DB, providers.Logger),
			providers.JWTProviderInstance,
		).Execute(c, locales.EN_US, token)

		if uc_result.HasError() {
			headers := map[api.HTTPHeaderTypeEnum]string{
				api.CONTENT_TYPE: string(api.APPLICATION_JSON),
			}
			content, statusCode := api.NewRequestResolver[models.UserWithRole]().ResolveDTO(c, uc_result, headers)
			c.JSON(statusCode, content)
			c.Abort()
			return
		}

		user := uc_result.GetData()

		ctx := context.Background()
		ctx = context.WithValue(ctx, app_context.UserKey, *user)

		c.Request = c.Request.WithContext(ctx)

		c.Next()

	}
}
