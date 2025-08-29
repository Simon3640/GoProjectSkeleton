package password_pipes

import (
	"context"

	usecases_password "gormgoskeleton/src/application/modules/password/use_cases"

	dtos "gormgoskeleton/src/application/shared/DTOs"
	"gormgoskeleton/src/application/shared/locales"
	usecase "gormgoskeleton/src/application/shared/use_case"
)

func NewGetResetPasswordPipe(
	ctx context.Context,
	locale locales.LocaleTypeEnum,
	create_password_token_uc *usecases_password.CreatePasswordTokenUseCase,
	create_password_uc *usecases_password.CreatePasswordUseCase,
) *usecase.DAG[dtos.PasswordTokenCreate, bool] {
	dag := usecase.NewDag(usecase.NewStep(create_password_token_uc), locale, ctx)
	return usecase.Then(dag, usecase.NewStep(create_password_uc))
}
