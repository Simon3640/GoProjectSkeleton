package middlewares

import (
	"context"

	api "github.com/simon3640/goprojectskeleton/gin"
	authusecases "github.com/simon3640/goprojectskeleton/src/application/modules/auth/use_cases"
	app_context "github.com/simon3640/goprojectskeleton/src/application/shared/context"
	"github.com/simon3640/goprojectskeleton/src/application/shared/locales"
	"github.com/simon3640/goprojectskeleton/src/application/shared/observability"
	usecase "github.com/simon3640/goprojectskeleton/src/application/shared/use_case"
	"github.com/simon3640/goprojectskeleton/src/domain/models"
	database "github.com/simon3640/goprojectskeleton/src/infrastructure/databases/goprojectskeleton"
	userrepositories "github.com/simon3640/goprojectskeleton/src/infrastructure/databases/goprojectskeleton/repositories/user"
	"github.com/simon3640/goprojectskeleton/src/infrastructure/providers"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("Authorization")

		appContext := app_context.AppContext{Context: c.Request.Context()}
		uc := authusecases.NewAuthUserUseCase(
			userrepositories.NewUserRepository(database.GoProjectSkeletondb.DB, providers.Logger),
			providers.JWTProviderInstance,
		)
		uc_result := usecase.InstrumentUseCase(
			uc,
			&appContext,
			locales.EN_US,
			token,
			observability.GetObservabilityComponents().Tracer,
			observability.GetObservabilityComponents().Metrics,
			observability.GetObservabilityComponents().Clock,
			"auth_user_use_case",
		)

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
