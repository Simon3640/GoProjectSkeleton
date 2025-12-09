package password_pipes

import (
	"context"

	usecases_password "github.com/simon3640/goprojectskeleton/src/application/modules/password/use_cases"

	dtos "github.com/simon3640/goprojectskeleton/src/application/shared/DTOs"
	"github.com/simon3640/goprojectskeleton/src/application/shared/locales"
	usecase "github.com/simon3640/goprojectskeleton/src/application/shared/use_case"
)

func NewGetResetPasswordPipe(
	ctx context.Context,
	locale locales.LocaleTypeEnum,
	create_password_token_uc *usecases_password.CreatePasswordTokenUseCase,
	create_password_uc *usecases_password.CreatePasswordUseCase,
) *usecase.DAG[dtos.PasswordTokenCreate, bool] {
	return nil
	// return nil --- IGNORE ---
	// dag := usecase.NewDag(usecase.NewStep(create_password_token_uc), locale, ctx)
	// return usecase.Then(dag, usecase.NewStep(create_password_uc))
}
