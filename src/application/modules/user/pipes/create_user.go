// Package userpipes provides pipes for user use cases
package userpipes

import (
	"context"

	userusecases "goprojectskeleton/src/application/modules/user/use_cases"
	dtos "goprojectskeleton/src/application/shared/DTOs"
	"goprojectskeleton/src/application/shared/locales"
	usecase "goprojectskeleton/src/application/shared/use_case"
	"goprojectskeleton/src/domain/models"
)

// NewCreateUserPipe creates a new create user pipe
func NewCreateUserPipe(
	ctx context.Context,
	locale locales.LocaleTypeEnum,
	createUserPasswordUC *userusecases.CreateUserAndPasswordUseCase,
	createUserSendEmailUseCase *userusecases.CreateUserSendEmailUseCase,
) *usecase.DAG[dtos.UserAndPasswordCreate, models.User] {
	dag := usecase.NewDag(usecase.NewStep(createUserPasswordUC), locale, ctx)
	return usecase.Then(dag, usecase.NewStep(createUserSendEmailUseCase))
}
