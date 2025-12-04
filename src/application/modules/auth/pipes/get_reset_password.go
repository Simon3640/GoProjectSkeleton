// Package authpipes provides pipes for the authentication module.
package authpipes

import (
	"context"

	authusecases "goprojectskeleton/src/application/modules/auth/use_cases"

	"goprojectskeleton/src/application/shared/locales"
	usecase "goprojectskeleton/src/application/shared/use_case"
)

// NewGetResetPasswordPipe creates a new get reset password pipe.
func NewGetResetPasswordPipe(
	ctx context.Context,
	locale locales.LocaleTypeEnum,
	getResetPasswordTokenUC *authusecases.GetResetPasswordTokenUseCase,
	getResetPasswordTokenSendEmailUC *authusecases.GetResetPasswordSendEmailUseCase,
) *usecase.DAG[string, bool] {
	dag := usecase.NewDag(usecase.NewStep(getResetPasswordTokenUC), locale, ctx)
	return usecase.Then(dag, usecase.NewStep(getResetPasswordTokenSendEmailUC))
}
