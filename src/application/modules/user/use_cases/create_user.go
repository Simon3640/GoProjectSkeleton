package userusecases

import (
	"context"
	"strings"

	contractsProviders "goprojectskeleton/src/application/contracts/providers"
	contracts_repositories "goprojectskeleton/src/application/contracts/repositories"
	dtos "goprojectskeleton/src/application/shared/DTOs"
	"goprojectskeleton/src/application/shared/locales"
	"goprojectskeleton/src/application/shared/locales/messages"
	"goprojectskeleton/src/application/shared/status"
	usecase "goprojectskeleton/src/application/shared/use_case"
	"goprojectskeleton/src/domain/models"
)

// CreateUserUseCase is a use case that creates a user
type CreateUserUseCase struct {
	appMessages *locales.Locale
	log         contractsProviders.ILoggerProvider
	repo        contracts_repositories.IUserRepository
	locale      locales.LocaleTypeEnum
}

var _ usecase.BaseUseCase[dtos.UserCreate, models.User] = (*CreateUserUseCase)(nil)

// SetLocale sets the locale for the use case
func (uc *CreateUserUseCase) SetLocale(locale locales.LocaleTypeEnum) {
	if locale != "" {
		uc.locale = locale
	}
}

// Execute executes the use case
func (uc *CreateUserUseCase) Execute(ctx context.Context,
	locale locales.LocaleTypeEnum,
	input dtos.UserCreate,
) *usecase.UseCaseResult[models.User] {
	result := usecase.NewUseCaseResult[models.User]()
	uc.SetLocale(locale)
	uc.validate(input, result)

	if result.HasError() {
		return result
	}

	res, err := uc.repo.Create(input)

	if err != nil {
		uc.log.Error("Error creating user", err.ToError())
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

func (uc *CreateUserUseCase) validate(
	input dtos.UserCreate,
	result *usecase.UseCaseResult[models.User]) {
	msgs := input.Validate()

	if len(msgs) > 0 {
		result.SetError(
			status.InvalidInput,
			strings.Join(msgs, "\n"),
		)
	}
}

// NewCreateUserUseCase creates a new create user use case
func NewCreateUserUseCase(
	log contractsProviders.ILoggerProvider,
	repo contracts_repositories.IUserRepository,
) *CreateUserUseCase {
	return &CreateUserUseCase{
		appMessages: locales.NewLocale(locales.EN_US),
		log:         log,
		repo:        repo,
	}
}
