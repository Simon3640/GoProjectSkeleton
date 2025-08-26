package usecases_password

import (
	"context"

	contracts_providers "gormgoskeleton/src/application/contracts/providers"
	contracts_repositories "gormgoskeleton/src/application/contracts/repositories"
	"gormgoskeleton/src/application/shared/guards"
	"gormgoskeleton/src/application/shared/locales"
	"gormgoskeleton/src/application/shared/locales/messages"
	"gormgoskeleton/src/application/shared/status"
	usecase "gormgoskeleton/src/application/shared/use_case"
	"gormgoskeleton/src/domain/models"
)

type CreatePasswordUseCase struct {
	usecase.BaseUseCaseValidation[models.PasswordCreateNoHash, bool]
	log          contracts_providers.ILoggerProvider
	repo         contracts_repositories.IPasswordRepository
	hashProvider contracts_providers.IHashProvider
}

var _ usecase.BaseUseCase[models.PasswordCreateNoHash, bool] = (*CreatePasswordUseCase)(nil)

func (uc *CreatePasswordUseCase) SetLocale(locale locales.LocaleTypeEnum) {
	if locale != "" {
		uc.Locale = locale
	}
}

func (uc *CreatePasswordUseCase) Execute(ctx context.Context,
	locale locales.LocaleTypeEnum,
	input models.PasswordCreateNoHash,
) *usecase.UseCaseResult[bool] {
	result := usecase.NewUseCaseResult[bool]()
	uc.SetLocale(locale)
	uc.Validate(ctx, input, result)
	if result.HasError() {
		return result
	}

	hashedPassword, err := uc.hashProvider.HashPassword(input.NoHashedPassword)
	if err != nil {
		uc.log.Error("Error hashing password", err.ToError())
		result.SetError(
			err.Code,
			uc.AppMessages.Get(
				uc.Locale,
				err.Context,
			),
		)
		return result
	}

	passwordCreate := models.NewPasswordCreate(
		input.UserID,
		hashedPassword,
		input.ExpiresAt,
		input.IsActive,
	)

	res, err := uc.repo.Create(passwordCreate)

	if err != nil {
		result.SetError(
			err.Code,
			uc.AppMessages.Get(
				uc.Locale,
				err.Context,
			),
		)
		return result
	}

	var success bool
	if res != nil {
		success = true
	} else {
		success = false
	}

	result.SetData(
		status.Success,
		success,
		uc.AppMessages.Get(
			uc.Locale,
			messages.MessageKeysInstance.PASSWORD_CREATED,
		),
	)
	return result
}

func NewCreatePasswordUseCase(
	log contracts_providers.ILoggerProvider,
	repo contracts_repositories.IPasswordRepository,
	hashProvider contracts_providers.IHashProvider,
) *CreatePasswordUseCase {
	return &CreatePasswordUseCase{
		BaseUseCaseValidation: usecase.BaseUseCaseValidation[models.PasswordCreateNoHash, bool]{
			AppMessages: locales.NewLocale(locales.EN_US),
			Guards: usecase.NewGuards(
				guards.RoleGuard("admin", "user"),
				guards.UserResourceGuard[models.PasswordCreateNoHash](),
			),
		},
		log:          log,
		repo:         repo,
		hashProvider: hashProvider,
	}
}
