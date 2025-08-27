package usecases_user

import (
	"context"

	contracts_providers "gormgoskeleton/src/application/contracts/providers"
	contracts_repositories "gormgoskeleton/src/application/contracts/repositories"
	"gormgoskeleton/src/application/shared/guards"
	"gormgoskeleton/src/application/shared/locales"
	"gormgoskeleton/src/application/shared/locales/messages"
	"gormgoskeleton/src/application/shared/status"
	usecase "gormgoskeleton/src/application/shared/use_case"
	"gormgoskeleton/src/domain/models"
)

type GetAllUserUseCase struct {
	usecase.BaseUseCaseValidation[Nil, []models.User]
	log  contracts_providers.ILoggerProvider
	repo contracts_repositories.IUserRepository
}

type Nil struct{}

var _ usecase.BaseUseCase[Nil, []models.User] = (*GetAllUserUseCase)(nil)

func (uc *GetAllUserUseCase) SetLocale(locale locales.LocaleTypeEnum) {
	if locale != "" {
		uc.Locale = locale
	}
}

func (uc *GetAllUserUseCase) Execute(
	ctx context.Context,
	locale locales.LocaleTypeEnum,
	input Nil,
) *usecase.UseCaseResult[[]models.User] {
	result := usecase.NewUseCaseResult[[]models.User]()
	uc.SetLocale(locale)
	uc.Validate(ctx, input, result)
	if result.HasError() {
		return result
	}

	data, err := uc.repo.GetAll(nil, nil, nil)
	if err != nil {
		uc.log.Error("Error getting all users", err.ToError())
		result.SetError(
			err.Code,
			uc.AppMessages.Get(
				uc.Locale,
				err.Context,
			),
		)
		return result
	}
	result.SetData(
		status.Success,
		data,
		uc.AppMessages.Get(
			uc.Locale,
			messages.MessageKeysInstance.USER_LIST_SUCCESS,
		),
	)
	return result
}

func NewGetAllUserUseCase(
	log contracts_providers.ILoggerProvider,
	repo contracts_repositories.IUserRepository,
) *GetAllUserUseCase {
	return &GetAllUserUseCase{
		BaseUseCaseValidation: usecase.BaseUseCaseValidation[Nil, []models.User]{
			AppMessages: locales.NewLocale(locales.EN_US),
			Guards:      usecase.NewGuards(guards.RoleGuard("admin")),
		},
		log:  log,
		repo: repo,
	}
}
