package usecases_user

import (
	"context"
	"strings"

	"gormgoskeleton/src/application/contracts"
	contracts_repositories "gormgoskeleton/src/application/contracts/repositories"
	"gormgoskeleton/src/application/shared/locales"
	"gormgoskeleton/src/application/shared/locales/messages"
	"gormgoskeleton/src/application/shared/status"
	usecase "gormgoskeleton/src/application/shared/use_case"
	"gormgoskeleton/src/domain/models"
)

type CreateUserAndPasswordUseCase struct {
	appMessages  *locales.Locale
	log          contracts.ILoggerProvider
	repo         contracts_repositories.IUserRepository
	hashProvider contracts.IHashProvider
	locale       locales.LocaleTypeEnum
}

var _ usecase.BaseUseCase[models.UserAndPasswordCreate, models.User] = (*CreateUserAndPasswordUseCase)(nil)

func (uc *CreateUserAndPasswordUseCase) SetLocale(locale locales.LocaleTypeEnum) {
	if locale != "" {
		uc.locale = locale
	}
}

func (uc *CreateUserAndPasswordUseCase) Execute(ctx context.Context,
	locale locales.LocaleTypeEnum,
	input models.UserAndPasswordCreate,
) *usecase.UseCaseResult[models.User] {
	result := usecase.NewUseCaseResult[models.User]()
	uc.SetLocale(locale)
	uc.validate(input, result)

	if result.HasError() {
		return result
	}

	input.Password, _ = uc.hashProvider.HashPassword(input.Password)

	res, err := uc.repo.CreateWithPassword(input)

	if err != nil {
		uc.log.Error("Error creating user with password", err.ToError())
		result.SetError(
			err.Code,
			uc.appMessages.Get(
				uc.locale,
				err.Context,
			),
		)
		return result
	}
	result.SetData(
		status.Success,
		*res,
		uc.appMessages.Get(
			uc.locale,
			messages.MessageKeysInstance.USER_WAS_CREATED,
		),
	)
	return result
}

func (uc *CreateUserAndPasswordUseCase) validate(
	input models.UserAndPasswordCreate,
	result *usecase.UseCaseResult[models.User]) {
	msgs := input.Validate()

	if len(msgs) > 0 {
		result.SetError(
			status.InvalidInput,
			strings.Join(msgs, "\n"),
		)
	}
}

func NewCreateUserAndPasswordUseCase(
	log contracts.ILoggerProvider,
	repo contracts_repositories.IUserRepository,
	hashProvider contracts.IHashProvider,
) *CreateUserAndPasswordUseCase {
	return &CreateUserAndPasswordUseCase{
		appMessages:  locales.NewLocale(locales.EN_US),
		log:          log,
		repo:         repo,
		hashProvider: hashProvider,
	}
}
