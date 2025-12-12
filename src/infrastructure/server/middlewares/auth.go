package middlewares

import (
	"context"

	api "github.com/simon3640/goprojectskeleton/gin"
	authusecases "github.com/simon3640/goprojectskeleton/src/application/modules/auth/use_cases"
	app_context "github.com/simon3640/goprojectskeleton/src/application/shared/context"
	"github.com/simon3640/goprojectskeleton/src/application/shared/locales"
	"github.com/simon3640/goprojectskeleton/src/domain/models"
	database "github.com/simon3640/goprojectskeleton/src/infrastructure/databases/goprojectskeleton"
	repositories "github.com/simon3640/goprojectskeleton/src/infrastructure/databases/goprojectskeleton/repositories"
	"github.com/simon3640/goprojectskeleton/src/infrastructure/providers"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("Authorization")

		uc_result := authusecases.NewAuthUserUseCase(
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
