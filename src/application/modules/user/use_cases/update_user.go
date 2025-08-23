package usecases_user

import (
	"context"

	"gormgoskeleton/src/application/contracts"
	contracts_repositories "gormgoskeleton/src/application/contracts/repositories"
	"gormgoskeleton/src/application/shared/guards"
	"gormgoskeleton/src/application/shared/locales"
	"gormgoskeleton/src/application/shared/locales/messages"
	"gormgoskeleton/src/application/shared/status"
	usecase "gormgoskeleton/src/application/shared/use_case"
	"gormgoskeleton/src/domain/models"
)

type UpdateUserUseCase struct {
	usecase.BaseUseCaseValidation[models.UserUpdate, models.User]
	log  contracts.ILoggerProvider
	repo contracts_repositories.IUserRepository
}

var _ usecase.BaseUseCase[models.UserUpdate, models.User] = (*UpdateUserUseCase)(nil)

func (uc *UpdateUserUseCase) SetLocale(locale locales.LocaleTypeEnum) {
	if locale != "" {
		uc.Locale = locale
	}
}

func (uc *UpdateUserUseCase) Execute(ctx context.Context,
	locale locales.LocaleTypeEnum,
	input models.UserUpdate,
) *usecase.UseCaseResult[models.User] {
	result := usecase.NewUseCaseResult[models.User]()
	uc.SetLocale(locale)
	uc.Validate(ctx, input, result)
	if result.HasError() {
		return result
	}

	res, err := uc.repo.Update(input.ID, input)

	if err != nil {
		uc.log.Error("Error updating user", err.ToError())
		result.SetError(
			err.Code,
			uc.AppMessages.Get(
				uc.Locale,
				err.Context,
			),
		)
	}
	result.SetData(
		status.Success,
		*res,
		uc.AppMessages.Get(uc.Locale, messages.MessageKeysInstance.USER_WAS_CREATED))
	return result
}

func NewUpdateUserUseCase(
	log contracts.ILoggerProvider,
	repo contracts_repositories.IUserRepository,
) *UpdateUserUseCase {
	return &UpdateUserUseCase{
		BaseUseCaseValidation: usecase.BaseUseCaseValidation[models.UserUpdate, models.User]{
			AppMessages: locales.NewLocale(locales.EN_US),
			Guards:      usecase.NewGuards(guards.RoleGuard("admin", "user"), guards.UserResourceGuard[models.UserUpdate]()),
		},
		log:  log,
		repo: repo,
	}
}
