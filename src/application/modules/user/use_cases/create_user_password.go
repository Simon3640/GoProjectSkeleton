package userusecases

import (
	"context"
	"strings"

	contractsProviders "github.com/simon3640/goprojectskeleton/src/application/contracts/providers"
	contracts_repositories "github.com/simon3640/goprojectskeleton/src/application/contracts/repositories"
	dtos "github.com/simon3640/goprojectskeleton/src/application/shared/DTOs"
	"github.com/simon3640/goprojectskeleton/src/application/shared/locales"
	"github.com/simon3640/goprojectskeleton/src/application/shared/locales/messages"
	"github.com/simon3640/goprojectskeleton/src/application/shared/status"
	usecase "github.com/simon3640/goprojectskeleton/src/application/shared/use_case"
	"github.com/simon3640/goprojectskeleton/src/domain/models"
)

// CreateUserAndPasswordUseCase is a use case that creates a user and a password
type CreateUserAndPasswordUseCase struct {
	appMessages  *locales.Locale
	log          contractsProviders.ILoggerProvider
	repo         contracts_repositories.IUserRepository
	hashProvider contractsProviders.IHashProvider
	locale       locales.LocaleTypeEnum
}

var _ usecase.BaseUseCase[dtos.UserAndPasswordCreate, models.User] = (*CreateUserAndPasswordUseCase)(nil)

// SetLocale sets the locale for the use case
func (uc *CreateUserAndPasswordUseCase) SetLocale(locale locales.LocaleTypeEnum) {
	if locale != "" {
		uc.locale = locale
	}
}

// Execute executes the use case
func (uc *CreateUserAndPasswordUseCase) Execute(ctx context.Context,
	locale locales.LocaleTypeEnum,
	input dtos.UserAndPasswordCreate,
) *usecase.UseCaseResult[models.User] {
	result := usecase.NewUseCaseResult[models.User]()
	uc.SetLocale(locale)
	uc.validate(&input, result)

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
	input *dtos.UserAndPasswordCreate,
	result *usecase.UseCaseResult[models.User]) {
	msgs := input.Validate()

	if len(msgs) > 0 {
		result.SetError(
			status.InvalidInput,
			strings.Join(msgs, "\n"),
		)
	}
}

// NewCreateUserAndPasswordUseCase creates a new create user and password use case
func NewCreateUserAndPasswordUseCase(
	log contractsProviders.ILoggerProvider,
	repo contracts_repositories.IUserRepository,
	hashProvider contractsProviders.IHashProvider,
) *CreateUserAndPasswordUseCase {
	return &CreateUserAndPasswordUseCase{
		appMessages:  locales.NewLocale(locales.EN_US),
		log:          log,
		repo:         repo,
		hashProvider: hashProvider,
	}
}
