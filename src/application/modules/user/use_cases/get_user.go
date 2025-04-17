package usecases_user

import (
	"context"
	"gormgoskeleton/src/application/contracts"
	contracts_repositories "gormgoskeleton/src/application/contracts/repositories"
	"gormgoskeleton/src/application/shared/locales"
	"gormgoskeleton/src/application/shared/locales/messages"
	"gormgoskeleton/src/application/shared/status"
	usecase "gormgoskeleton/src/application/shared/use_case"
	"gormgoskeleton/src/domain/models"
)

type GetUserUseCase struct {
	appMessages *locales.Locale
	log         contracts.ILoggerProvider
	repo        contracts_repositories.IUserRepository[models.UserCreate, models.UserUpdate, models.User, any]
	locale      locales.LocaleTypeEnum
}

func (uc *GetUserUseCase) SetLocale(locale locales.LocaleTypeEnum) {
	if locale != "" {
		uc.locale = locale
	}
}

func (uc *GetUserUseCase) Execute(ctx context.Context,
	locale locales.LocaleTypeEnum,
	input int,
) *usecase.UseCaseResult[models.User] {
	result := usecase.NewUseCaseResult[models.User]()
	uc.GetUser(ctx, result, input)
	uc.SetLocale(locale)
	return result
}

func (uc *GetUserUseCase) GetUser(ctx context.Context, result *usecase.UseCaseResult[models.User], id int) {
	res, err := uc.repo.GetByID(id)
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
		*res,
		"",
	)
}

func NewGetUserUseCase(
	log contracts.ILoggerProvider,
	repo contracts_repositories.IUserRepository[models.UserCreate, models.UserUpdate, models.User, any],
) *GetUserUseCase {
	return &GetUserUseCase{
		appMessages: locales.NewLocale(locales.EN_US),
		log:         log,
		repo:        repo,
	}
}
