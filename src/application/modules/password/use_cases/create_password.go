package passwordusecases

import (
	contractsproviders "github.com/simon3640/goprojectskeleton/src/application/contracts/providers"
	passwordcontracts "github.com/simon3640/goprojectskeleton/src/application/modules/password/contracts"
	dtos "github.com/simon3640/goprojectskeleton/src/application/modules/password/dtos"
	passwordservices "github.com/simon3640/goprojectskeleton/src/application/modules/password/services"
	app_context "github.com/simon3640/goprojectskeleton/src/application/shared/context"
	"github.com/simon3640/goprojectskeleton/src/application/shared/guards"
	"github.com/simon3640/goprojectskeleton/src/application/shared/locales"
	"github.com/simon3640/goprojectskeleton/src/application/shared/locales/messages"
	"github.com/simon3640/goprojectskeleton/src/application/shared/status"
	usecase "github.com/simon3640/goprojectskeleton/src/application/shared/use_case"
)

// CreatePasswordUseCase is the use case for creating a password
type CreatePasswordUseCase struct {
	usecase.BaseUseCaseValidation[dtos.PasswordCreateNoHash, bool]
	log          contractsproviders.ILoggerProvider
	repo         passwordcontracts.IPasswordRepository
	hashProvider contractsproviders.IHashProvider
}

var _ usecase.BaseUseCase[dtos.PasswordCreateNoHash, bool] = (*CreatePasswordUseCase)(nil)

func (uc *CreatePasswordUseCase) SetLocale(locale locales.LocaleTypeEnum) {
	if locale != "" {
		uc.Locale = locale
	}
}

// Execute executes the use case
func (uc *CreatePasswordUseCase) Execute(ctx *app_context.AppContext,
	locale locales.LocaleTypeEnum,
	input dtos.PasswordCreateNoHash,
) *usecase.UseCaseResult[bool] {
	result := usecase.NewUseCaseResult[bool]()
	uc.SetLocale(locale)
	uc.Validate(ctx, input, result)
	if result.HasError() {
		return result
	}

	uc.createPassword(input, result)
	if result.HasError() {
		return result
	}

	uc.setSuccessResult(result)
	return result
}

func (uc *CreatePasswordUseCase) createPassword(input dtos.PasswordCreateNoHash, result *usecase.UseCaseResult[bool]) {
	_, err := passwordservices.CreatePasswordService(input, uc.hashProvider, uc.repo)

	if err != nil {
		uc.log.Error("CreatePasswordUseCase: Execute: Error creating password", err.ToError())
		result.SetError(
			err.Code,
			uc.AppMessages.Get(
				uc.Locale,
				err.Context,
			),
		)
	}
}

func (uc *CreatePasswordUseCase) setSuccessResult(result *usecase.UseCaseResult[bool]) {
	result.SetData(
		status.Success,
		true,
		uc.AppMessages.Get(
			uc.Locale,
			messages.MessageKeysInstance.PASSWORD_CREATED,
		),
	)
}

func NewCreatePasswordUseCase(
	log contractsproviders.ILoggerProvider,
	repo passwordcontracts.IPasswordRepository,
	hashProvider contractsproviders.IHashProvider,
	skip_guards bool,
) *CreatePasswordUseCase {
	return &CreatePasswordUseCase{
		BaseUseCaseValidation: usecase.BaseUseCaseValidation[dtos.PasswordCreateNoHash, bool]{
			AppMessages: locales.NewLocale(locales.EN_US),
			Guards: usecase.NewGuards(
				guards.RoleGuard("admin", "user"),
				guards.UserResourceGuard[dtos.PasswordCreateNoHash](),
			),
		},
		log:          log,
		repo:         repo,
		hashProvider: hashProvider,
	}
}
