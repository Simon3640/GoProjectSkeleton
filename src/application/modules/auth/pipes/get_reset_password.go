// Package authpipes provides pipes for the authentication module.
package authpipes

import (
	"context"

	authusecases "github.com/simon3640/goprojectskeleton/src/application/modules/auth/use_cases"
	shareddtos "github.com/simon3640/goprojectskeleton/src/application/shared/DTOs"

	"github.com/simon3640/goprojectskeleton/src/application/shared/locales"
	usecase "github.com/simon3640/goprojectskeleton/src/application/shared/use_case"
)

// NewGetResetPasswordPipe creates a new get reset password pipe.
// El token se genera de forma síncrona, y el email se envía en background.
// Esto permite retornar la respuesta inmediatamente sin esperar el envío del email.
// Retorna el tipo OneTimeTokenUser (que incluye el token) en lugar de bool.
func NewGetResetPasswordPipe(
	ctx context.Context,
	locale locales.LocaleTypeEnum,
	getResetPasswordTokenUC *authusecases.GetResetPasswordTokenUseCase,
	getResetPasswordTokenSendEmailUC *authusecases.GetResetPasswordSendEmailUseCase,
) *usecase.DAG[string, shareddtos.OneTimeTokenUser] {
	dag := usecase.NewDag(usecase.NewStep(getResetPasswordTokenUC), locale, ctx)
	// El email se envía en background, la respuesta se retorna inmediatamente con el token
	return usecase.ThenBackground(dag, usecase.NewStep(getResetPasswordTokenSendEmailUC))
}
