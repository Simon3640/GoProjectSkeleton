package usecases_user

import (
	"context"
	"go/types"
	"gormgoskeleton/src/application/contracts"
	contracts_repositories "gormgoskeleton/src/application/contracts/repositories"
	"gormgoskeleton/src/application/shared/locales"
	"gormgoskeleton/src/application/shared/locales/messages"
	"gormgoskeleton/src/application/shared/status"
	usecase "gormgoskeleton/src/application/shared/use_case"
)

type DeleteUserUseCase struct {
	appMessages *locales.Locale
	log         contracts.ILoggerProvider
	repo        contracts_repositories.IUserRepository
	locale      locales.LocaleTypeEnum
}

var _ usecase.BaseUseCase[int, types.Nil] = (*DeleteUserUseCase)(nil)

func (uc *DeleteUserUseCase) SetLocale(locale locales.LocaleTypeEnum) {
	if locale != "" {
		uc.locale = locale
	}
}

func (uc *DeleteUserUseCase) Execute(ctx context.Context,
	locale locales.LocaleTypeEnum,
	input int,
) *usecase.UseCaseResult[types.Nil] {
	result := usecase.NewUseCaseResult[types.Nil]()
	uc.SetLocale(locale)
	err := uc.repo.Delete(input)
	if err != nil {
		result.SetError(
			status.Conflict,
			uc.appMessages.Get(
				uc.locale,
				messages.MessageKeysInstance.SOMETHING_WENT_WRONG,
			),
		)
	}
	result.SetData(
		status.Success,
		types.Nil{},
		uc.appMessages.Get(
			uc.locale,
			messages.MessageKeysInstance.USER_DELETE_SUCCESS,
		),
	)
	return result
}

func NewDeleteUserUseCase(
	log contracts.ILoggerProvider,
	repo contracts_repositories.IUserRepository,
) *DeleteUserUseCase {
	return &DeleteUserUseCase{
		appMessages: locales.NewLocale(locales.EN_US),
		log:         log,
		repo:        repo,
	}
}
