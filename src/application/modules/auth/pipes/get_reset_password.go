// Package authpipes provides pipes for the authentication module.
package authpipes

import (
	authusecases "github.com/simon3640/goprojectskeleton/src/application/modules/auth/use_cases"
	shareddtos "github.com/simon3640/goprojectskeleton/src/application/shared/DTOs"
	app_context "github.com/simon3640/goprojectskeleton/src/application/shared/context"

	"github.com/simon3640/goprojectskeleton/src/application/shared/locales"
	usecase "github.com/simon3640/goprojectskeleton/src/application/shared/use_case"
	workers "github.com/simon3640/goprojectskeleton/src/application/shared/workers"
)

// NewGetResetPasswordPipe creates a new get reset password pipe.
// El token se genera de forma síncrona, y el email se envía en background.
// Esto permite retornar la respuesta inmediatamente sin esperar el envío del email.
// Retorna el tipo OneTimeTokenUser (que incluye el token) en lugar de bool.
func NewGetResetPasswordPipe(
	ctx *app_context.AppContext,
	locale locales.LocaleTypeEnum,
	getResetPasswordTokenUC *authusecases.GetResetPasswordTokenUseCase,
	getResetPasswordTokenSendEmailUC *authusecases.GetResetPasswordSendEmailUseCase,
) *usecase.DAG[string, shareddtos.OneTimeTokenUser] {
	dag := usecase.NewDag(ctx, usecase.NewStep(getResetPasswordTokenUC), locale, workers.GetBackgroundExecutor())
	// The email is sent in background, the response is returned immediately with the token
	return usecase.ThenBackground(dag, usecase.NewStep(getResetPasswordTokenSendEmailUC), "get-reset-password-send-email")
}
