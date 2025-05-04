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

type GetAllUserUseCase struct {
	appMessages *locales.Locale
	log         contracts.ILoggerProvider
	repo        contracts_repositories.IUserRepository[models.UserCreate, models.UserUpdate, models.User, any]
	locale      locales.LocaleTypeEnum
}

type Nil struct{}

var _ usecase.BaseUseCase[Nil, []models.User] = (*GetAllUserUseCase)(nil)

func (uc *GetAllUserUseCase) SetLocale(locale locales.LocaleTypeEnum) {
	if locale != "" {
		uc.locale = locale
	}
}

func (uc *GetAllUserUseCase) Execute(
	ctx context.Context,
	locale locales.LocaleTypeEnum,
	input Nil,
) *usecase.UseCaseResult[[]models.User] {
	result := usecase.NewUseCaseResult[[]models.User]()
	uc.SetLocale(locale)
	data, err := uc.repo.GetAll(nil, nil, nil)
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
		data,
		uc.appMessages.Get(
			uc.locale,
			messages.MessageKeysInstance.USER_LIST_SUCCESS,
		),
	)
	return result
}

func NewGetAllUserUseCase(
	log contracts.ILoggerProvider,
	repo contracts_repositories.IUserRepository[models.UserCreate, models.UserUpdate, models.User, any],
) *GetAllUserUseCase {
	return &GetAllUserUseCase{
		log:         log,
		repo:        repo,
		appMessages: locales.NewLocale(locales.EN_US),
	}
}
