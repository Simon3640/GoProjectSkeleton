package userusecases

import (
	"strings"

	contractsproviders "github.com/simon3640/goprojectskeleton/src/application/contracts/providers"
	usercontracts "github.com/simon3640/goprojectskeleton/src/application/modules/user/contracts"
	userdtos "github.com/simon3640/goprojectskeleton/src/application/modules/user/dtos"
	app_context "github.com/simon3640/goprojectskeleton/src/application/shared/context"
	"github.com/simon3640/goprojectskeleton/src/application/shared/locales"
	"github.com/simon3640/goprojectskeleton/src/application/shared/locales/messages"
	"github.com/simon3640/goprojectskeleton/src/application/shared/status"
	usecase "github.com/simon3640/goprojectskeleton/src/application/shared/use_case"
	"github.com/simon3640/goprojectskeleton/src/domain/models"
)

// CreateUserUseCase is a use case that creates a user
type CreateUserUseCase struct {
	usecase.BaseUseCaseValidation[userdtos.UserCreate, models.User]
	log  contractsproviders.ILoggerProvider
	repo usercontracts.IUserRepository
}

var _ usecase.BaseUseCase[userdtos.UserCreate, models.User] = (*CreateUserUseCase)(nil)

// Execute executes the use case
func (uc *CreateUserUseCase) Execute(ctx *app_context.AppContext,
	locale locales.LocaleTypeEnum,
	input userdtos.UserCreate,
) *usecase.UseCaseResult[models.User] {
	result := usecase.NewUseCaseResult[models.User]()
	uc.SetLocale(locale)
	uc.SetAppContext(ctx)
	uc.validate(input, result)

	if result.HasError() {
		return result
	}

	res, err := uc.repo.Create(input)

	if err != nil {
		uc.log.Error("Error creating user", err.ToError())
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
		*res,
		uc.AppMessages.Get(
			uc.Locale,
			messages.MessageKeysInstance.USER_WAS_CREATED,
		),
	)
	return result
}

func (uc *CreateUserUseCase) validate(
	input userdtos.UserCreate,
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
	log contractsproviders.ILoggerProvider,
	repo usercontracts.IUserRepository,
) *CreateUserUseCase {
	return &CreateUserUseCase{
		BaseUseCaseValidation: usecase.BaseUseCaseValidation[userdtos.UserCreate, models.User]{
			AppMessages: locales.NewLocale(locales.EN_US),
			Guards:      usecase.NewGuards(),
		},
		log:  log,
		repo: repo,
	}
}
