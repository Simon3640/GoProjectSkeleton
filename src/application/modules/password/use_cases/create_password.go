package usecases_password

import (
	"context"

	contractsProviders "gormgoskeleton/src/application/contracts/providers"
	contracts_repositories "gormgoskeleton/src/application/contracts/repositories"
	dtos "gormgoskeleton/src/application/shared/DTOs"
	"gormgoskeleton/src/application/shared/guards"
	"gormgoskeleton/src/application/shared/locales"
	"gormgoskeleton/src/application/shared/locales/messages"
	"gormgoskeleton/src/application/shared/services"
	"gormgoskeleton/src/application/shared/status"
	usecase "gormgoskeleton/src/application/shared/use_case"
)

type CreatePasswordUseCase struct {
	usecase.BaseUseCaseValidation[dtos.PasswordCreateNoHash, bool]
	log          contractsProviders.ILoggerProvider
	repo         contracts_repositories.IPasswordRepository
	hashProvider contractsProviders.IHashProvider
}

var _ usecase.BaseUseCase[dtos.PasswordCreateNoHash, bool] = (*CreatePasswordUseCase)(nil)

func (uc *CreatePasswordUseCase) SetLocale(locale locales.LocaleTypeEnum) {
	if locale != "" {
		uc.Locale = locale
	}
}

func (uc *CreatePasswordUseCase) Execute(ctx context.Context,
	locale locales.LocaleTypeEnum,
	input dtos.PasswordCreateNoHash,
) *usecase.UseCaseResult[bool] {
	result := usecase.NewUseCaseResult[bool]()
	uc.SetLocale(locale)
	uc.Validate(ctx, input, result)
	if result.HasError() {
		return result
	}

	_, err := services.CreatePasswordService(input, uc.hashProvider, uc.repo)

	if err != nil {
		uc.log.Error("CreatePasswordUseCase: Execute: Error creating password", err.ToError())
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
		true,
		uc.AppMessages.Get(
			uc.Locale,
			messages.MessageKeysInstance.PASSWORD_CREATED,
		),
	)
	return result
}

func NewCreatePasswordUseCase(
	log contractsProviders.ILoggerProvider,
	repo contracts_repositories.IPasswordRepository,
	hashProvider contractsProviders.IHashProvider,
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
