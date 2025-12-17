package userusecases

import (
	"strings"

	contractsProviders "github.com/simon3640/goprojectskeleton/src/application/contracts/providers"
	usercontracts "github.com/simon3640/goprojectskeleton/src/application/modules/user/contracts"
	userdtos "github.com/simon3640/goprojectskeleton/src/application/modules/user/dtos"
	app_context "github.com/simon3640/goprojectskeleton/src/application/shared/context"
	applicationerror "github.com/simon3640/goprojectskeleton/src/application/shared/errors"
	"github.com/simon3640/goprojectskeleton/src/application/shared/locales"
	"github.com/simon3640/goprojectskeleton/src/application/shared/locales/messages"
	"github.com/simon3640/goprojectskeleton/src/application/shared/observability"
	"github.com/simon3640/goprojectskeleton/src/application/shared/status"
	usecase "github.com/simon3640/goprojectskeleton/src/application/shared/use_case"
	"github.com/simon3640/goprojectskeleton/src/domain/models"
)

// CreateUserAndPasswordUseCase is a use case that creates a user and a password
type CreateUserAndPasswordUseCase struct {
	usecase.BaseUseCaseValidation[userdtos.UserAndPasswordCreate, models.User]
	repo         usercontracts.IUserRepository
	hashProvider contractsProviders.IHashProvider
}

var _ usecase.BaseUseCase[userdtos.UserAndPasswordCreate, models.User] = (*CreateUserAndPasswordUseCase)(nil)

// Execute executes the use case
func (uc *CreateUserAndPasswordUseCase) Execute(ctx *app_context.AppContext,
	locale locales.LocaleTypeEnum,
	input userdtos.UserAndPasswordCreate,
) *usecase.UseCaseResult[models.User] {
	result := usecase.NewUseCaseResult[models.User]()
	uc.SetLocale(locale)
	uc.SetAppContext(ctx)
	uc.validate(&input, result)

	if result.HasError() {
		return result
	}

	uc.hashPassword(&input, result)
	if result.HasError() {
		return result
	}

	res := uc.createUser(input, result)
	if result.HasError() {
		return result
	}

	result.SetData(
		status.Success,
		*res,
		uc.AppMessages.Get(
			uc.Locale,
			messages.MessageKeysInstance.USER_WAS_CREATED,
		),
	)
	observability.GetObservabilityComponents().Logger.InfoWithContext("user_created_and_password_hashed", uc.AppContext)
	return result
}

func (uc *CreateUserAndPasswordUseCase) createUser(input userdtos.UserAndPasswordCreate, result *usecase.UseCaseResult[models.User]) *models.User {
	res, err := uc.repo.CreateWithPassword(input)
	if err != nil {
		observability.GetObservabilityComponents().Logger.ErrorWithContext("Error creating user with password", err.ToError(), uc.AppContext)
		result.SetError(
			err.Code,
			uc.AppMessages.Get(uc.Locale, err.Context),
		)
		return nil
	}
	return res
}

func (uc *CreateUserAndPasswordUseCase) hashPassword(input *userdtos.UserAndPasswordCreate, result *usecase.UseCaseResult[models.User]) {
	var err *applicationerror.ApplicationError
	input.Password, err = uc.hashProvider.HashPassword(input.Password)
	if err != nil {
		observability.GetObservabilityComponents().Logger.ErrorWithContext("Error hashing password", err.ToError(), uc.AppContext)
		result.SetError(
			err.Code,
			uc.AppMessages.Get(uc.Locale, err.Context),
		)
	}
}

func (uc *CreateUserAndPasswordUseCase) validate(
	input *userdtos.UserAndPasswordCreate,
	result *usecase.UseCaseResult[models.User]) {
	msgs := input.Validate()

	if len(msgs) > 0 {
		observability.GetObservabilityComponents().Logger.WarningWithContext("Invalid input", uc.AppContext)
		result.SetError(
			status.InvalidInput,
			strings.Join(msgs, "\n"),
		)
	}
}

// NewCreateUserAndPasswordUseCase creates a new create user and password use case
func NewCreateUserAndPasswordUseCase(
	repo usercontracts.IUserRepository,
	hashProvider contractsProviders.IHashProvider,
) *CreateUserAndPasswordUseCase {
	return &CreateUserAndPasswordUseCase{
		BaseUseCaseValidation: usecase.BaseUseCaseValidation[userdtos.UserAndPasswordCreate, models.User]{
			AppMessages: locales.NewLocale(locales.EN_US),
			Guards:      usecase.NewGuards(),
		},
		repo:         repo,
		hashProvider: hashProvider,
	}
}
