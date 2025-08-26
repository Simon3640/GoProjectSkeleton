package user_pipes

import (
	"context"

	usecases_user "gormgoskeleton/src/application/modules/user/use_cases"
	"gormgoskeleton/src/application/shared/locales"
	usecase "gormgoskeleton/src/application/shared/use_case"
	"gormgoskeleton/src/domain/models"
)

func NewCreateUserPipe(
	ctx context.Context,
	locale locales.LocaleTypeEnum,
	create_user_password_uc *usecases_user.CreateUserAndPasswordUseCase,
	create_user_send_email_uc *usecases_user.CreateUserSendEmailUseCase,
) *usecase.DAG[models.UserAndPasswordCreate, models.User] {
	dag := usecase.NewDag(usecase.NewStep(create_user_password_uc), locale, ctx)
	return usecase.Then(dag, usecase.NewStep(create_user_send_email_uc))
}
