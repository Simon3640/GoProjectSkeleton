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
	"strings"
)

type UpdateUserUseCase struct {
	appMessages *locales.Locale
	log         contracts.ILoggerProvider
	repo        contracts_repositories.IUserRepository
	locale      locales.LocaleTypeEnum
}

var _ usecase.BaseUseCase[models.UserUpdate, models.User] = (*UpdateUserUseCase)(nil)

func (uc *UpdateUserUseCase) SetLocale(locale locales.LocaleTypeEnum) {
	if locale != "" {
		uc.locale = locale
	}
}

func (uc *UpdateUserUseCase) Execute(ctx context.Context,
	locale locales.LocaleTypeEnum,
	input models.UserUpdate,
) *usecase.UseCaseResult[models.User] {
	result := usecase.NewUseCaseResult[models.User]()
	uc.SetLocale(locale)
	validation, msg := uc.validate(input)

	if !validation {
		result.SetError(
			status.InvalidInput,
			strings.Join(msg, "\n"),
		)
		return result
	}

	res, err := uc.repo.Update(input.ID, input)

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
		uc.appMessages.Get(uc.locale, messages.MessageKeysInstance.USER_WAS_CREATED))
	return result
}

func (uc *UpdateUserUseCase) validate(input models.UserUpdate) (bool, []string) {
	var msg []string
	if input.ID == 0 {
		msg = append(msg, uc.appMessages.Get(
			uc.locale,
			messages.MessageKeysInstance.SOME_PARAMETERS_ARE_MISSING,
		))
	}
	return len(msg) == 0, msg
}

func NewUpdateUserUseCase(
	log contracts.ILoggerProvider,
	repo contracts_repositories.IUserRepository,
) *UpdateUserUseCase {
	return &UpdateUserUseCase{
		appMessages: locales.NewLocale(locales.EN_US),
		log:         log,
		repo:        repo,
	}
}
