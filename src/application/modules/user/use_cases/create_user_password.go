package usecases_user

import (
	"context"
	"regexp"
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
	validation, msg := uc.validate(input)

	if !validation {
		result.SetError(
			status.InvalidInput,
			strings.Join(msg, "\n"),
		)
		return result
	}

	input.Password, _ = uc.hashProvider.HashPassword(input.Password)

	res, err := uc.repo.CreateWithPassword(input)

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
		uc.appMessages.Get(
			uc.locale,
			messages.MessageKeysInstance.USER_WAS_CREATED,
		),
	)
	return result
}

func (uc *CreateUserAndPasswordUseCase) validate(input models.UserAndPasswordCreate) (bool, []string) {
	var msgs []string

	if input.Email == "" {
		msgs = append(msgs, uc.appMessages.Get(uc.locale, messages.MessageKeysInstance.SOME_PARAMETERS_ARE_MISSING))
	}
	// regex for email validation
	if !regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`).MatchString(input.Email) {
		msgs = append(msgs, uc.appMessages.Get(uc.locale, messages.MessageKeysInstance.INVALID_EMAIL))
	}

	if !models.IsValidPassword(input.Password) {
		msgs = append(msgs, uc.appMessages.Get(uc.locale, messages.MessageKeysInstance.INVALID_PASSWORD))
	}

	return len(msgs) == 0, msgs
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
