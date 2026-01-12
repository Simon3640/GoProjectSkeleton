package userusecases

import (
	"strings"

	usercontracts "github.com/simon3640/goprojectskeleton/src/application/modules/user/contracts"
	userdtos "github.com/simon3640/goprojectskeleton/src/application/modules/user/dtos"
	app_context "github.com/simon3640/goprojectskeleton/src/application/shared/context"
	"github.com/simon3640/goprojectskeleton/src/application/shared/locales"
	"github.com/simon3640/goprojectskeleton/src/application/shared/locales/messages"
	"github.com/simon3640/goprojectskeleton/src/application/shared/observability"
	"github.com/simon3640/goprojectskeleton/src/application/shared/status"
	usecase "github.com/simon3640/goprojectskeleton/src/application/shared/use_case"
	usermodels "github.com/simon3640/goprojectskeleton/src/domain/user/models"
)

// CreateUserUseCase is a use case that creates a user
type CreateUserUseCase struct {
	usecase.BaseUseCaseValidation[userdtos.UserCreate, usermodels.User]
	repo usercontracts.IUserRepository
}

var _ usecase.BaseUseCase[userdtos.UserCreate, usermodels.User] = (*CreateUserUseCase)(nil)

// Execute executes the use case
func (uc *CreateUserUseCase) Execute(ctx *app_context.AppContext,
	locale locales.LocaleTypeEnum,
	input userdtos.UserCreate,
) *usecase.UseCaseResult[usermodels.User] {
	result := usecase.NewUseCaseResult[usermodels.User]()
	uc.SetLocale(locale)
	uc.SetAppContext(ctx)
	uc.validate(input, result)

	if result.HasError() {
		return result
	}

	res := uc.createUser(input, result)
	if result.HasError() {
		return result
	}

	result.SetData(
		status.Created,
		*res,
		uc.AppMessages.Get(
			uc.Locale,
			messages.MessageKeysInstance.USER_WAS_CREATED,
		),
	)
	return result
}

func (uc *CreateUserUseCase) createUser(input userdtos.UserCreate, result *usecase.UseCaseResult[usermodels.User]) *usermodels.User {
	res, err := uc.repo.Create(input)
	if err != nil {
		observability.GetObservabilityComponents().Logger.ErrorWithContext("Error creating user", err.ToError(), uc.AppContext)
		result.SetError(err.Code, uc.AppMessages.Get(uc.Locale, err.Context))
	}
	return res
}

func (uc *CreateUserUseCase) validate(
	input userdtos.UserCreate,
	result *usecase.UseCaseResult[usermodels.User]) {
	msgs := input.Validate()

	if len(msgs) > 0 {
		observability.GetObservabilityComponents().Logger.WarningWithContext("Invalid input", uc.AppContext)
		result.SetError(
			status.InvalidInput,
			strings.Join(msgs, "\n"),
		)
	}
}

// NewCreateUserUseCase creates a new create user use case
func NewCreateUserUseCase(
	repo usercontracts.IUserRepository,
) *CreateUserUseCase {
	return &CreateUserUseCase{
		BaseUseCaseValidation: usecase.BaseUseCaseValidation[userdtos.UserCreate, usermodels.User]{
			AppMessages: locales.NewLocale(locales.EN_US),
			Guards:      usecase.NewGuards(),
		},
		repo: repo,
	}
}
