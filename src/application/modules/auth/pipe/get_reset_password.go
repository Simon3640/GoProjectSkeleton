package auth_pipes

import (
	"context"

	"gormgoskeleton/src/application/modules/auth"

	"gormgoskeleton/src/application/shared/locales"
	usecase "gormgoskeleton/src/application/shared/use_case"
)

func NewGetResetPasswordPipe(
	ctx context.Context,
	locale locales.LocaleTypeEnum,
	get_reset_password_token_uc *auth.GetResetPasswordTokenUseCase,
	get_reset_password_token_send_email_uc *auth.GetResetPasswordSendEmailUseCase,
) *usecase.DAG[string, bool] {
	dag := usecase.NewDag(usecase.NewStep(get_reset_password_token_uc), locale, ctx)
	return usecase.Then(dag, usecase.NewStep(get_reset_password_token_send_email_uc))
}
