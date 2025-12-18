package authhandlers

import (
	"encoding/json"
	"net/http"

	authdtos "github.com/simon3640/goprojectskeleton/src/application/modules/auth/dtos"
	authusecases "github.com/simon3640/goprojectskeleton/src/application/modules/auth/use_cases"
	"github.com/simon3640/goprojectskeleton/src/application/shared/observability"
	usecase "github.com/simon3640/goprojectskeleton/src/application/shared/use_case"
	handlers "github.com/simon3640/goprojectskeleton/src/infrastructure/handlers/shared"
	"github.com/simon3640/goprojectskeleton/src/infrastructure/providers"
)

// RefreshAccessToken refresh JWT access token
// @Summary      Refresh JWT access token
// @Description  This endpoint allows a user to refresh their JWT access token using a valid refresh
// @Tags         Auth
// @Accept       json
// @Param Accept-Language header string false "Locale for response messages" Enums(en-US, es-ES) default(en-US)
// @Produce      json
// @Param        request body string true "Refresh token"
// @Success      200 {object} authdtos.Token
// @Failure      400 {object} map[string]string "Validation error"
// @Router       /api/auth/refresh [post]
func RefreshAccessToken(ctx handlers.HandlerContext) {
	var refreshToken string
	if err := json.NewDecoder(*ctx.Body).Decode(&refreshToken); err != nil {
		http.Error(ctx.ResponseWriter, err.Error(), http.StatusBadRequest)
		return
	}

	uc := authusecases.NewAuthenticationRefreshUseCase(
		providers.JWTProviderInstance,
	)
	ucResult := usecase.InstrumentUseCase(
		uc,
		ctx.Context,
		ctx.Locale,
		refreshToken,
		observability.GetObservabilityComponents().Tracer,
		observability.GetObservabilityComponents().Metrics,
		observability.GetObservabilityComponents().Clock,
		"authenticate_refresh_use_case",
	)
	headers := map[handlers.HTTPHeaderTypeEnum]string{
		handlers.CONTENT_TYPE: string(handlers.APPLICATION_JSON),
	}
	handlers.NewRequestResolver[authdtos.Token]().ResolveDTO(ctx.ResponseWriter, ucResult, headers)
}
