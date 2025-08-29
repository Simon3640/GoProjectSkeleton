package usecases_password

import (
	"context"
	"time"

	contracts_providers "gormgoskeleton/src/application/contracts/providers"
	contracts_repositories "gormgoskeleton/src/application/contracts/repositories"
	dtos "gormgoskeleton/src/application/shared/DTOs"
	"gormgoskeleton/src/application/shared/locales"
	"gormgoskeleton/src/application/shared/locales/messages"
	"gormgoskeleton/src/application/shared/status"
	usecase "gormgoskeleton/src/application/shared/use_case"
)

type CreatePasswordTokenUseCase struct {
	usecase.BaseUseCaseValidation[dtos.PasswordTokenCreate, dtos.PasswordCreateNoHash]
	log              contracts_providers.ILoggerProvider
	passRepo         contracts_repositories.IPasswordRepository
	hashProvider     contracts_providers.IHashProvider
	oneTimetokenRepo contracts_repositories.IOneTimeTokenRepository
}

var _ usecase.BaseUseCase[dtos.PasswordTokenCreate, dtos.PasswordCreateNoHash] = (*CreatePasswordTokenUseCase)(nil)

func (uc *CreatePasswordTokenUseCase) SetLocale(locale locales.LocaleTypeEnum) {
	if locale != "" {
		uc.Locale = locale
	}
}

func (uc *CreatePasswordTokenUseCase) Execute(ctx context.Context,
	locale locales.LocaleTypeEnum,
	input dtos.PasswordTokenCreate,
) *usecase.UseCaseResult[dtos.PasswordCreateNoHash] {
	result := usecase.NewUseCaseResult[dtos.PasswordCreateNoHash]()
	uc.SetLocale(locale)
	uc.Validate(ctx, input, result)
	if result.HasError() {
		return result
	}

	hash := uc.hashProvider.HashOneTimeToken(input.Token)
	oneTimeToke, err := uc.oneTimetokenRepo.GetByTokenHash(hash)
	if err != nil {
		uc.log.Error("Error getting one time token by hash", err.ToError())
		result.SetError(
			err.Code,
			uc.AppMessages.Get(
				uc.Locale,
				err.Context,
			),
		)
		return result
	}

	if oneTimeToke == nil || oneTimeToke.IsUsed || oneTimeToke.Expires.Before(time.Now()) {
		uc.log.Error("One time token is not valid", nil)
		result.SetError(
			status.Conflict,
			uc.AppMessages.Get(
				uc.Locale,
				messages.MessageKeysInstance.INVALID_PASSWORD_RESET_TOKEN,
			),
		)
		return result
	}

	passwordCreateNoHash := dtos.PasswordCreateNoHash{
		UserID:           oneTimeToke.UserID,
		NoHashedPassword: input.NoHashedPassword,
		IsActive:         true,
	}

	result.SetData(
		status.Success,
		passwordCreateNoHash,
		uc.AppMessages.Get(
			uc.Locale,
			messages.MessageKeysInstance.RESET_PASSWORD_TOKEN_VALID,
		),
	)

	return result
}

func (uc *CreatePasswordTokenUseCase) Validate(
	ctx context.Context,
	input dtos.PasswordTokenCreate,
	result *usecase.UseCaseResult[dtos.PasswordCreateNoHash],
) {
	// Call base validation
}

func NewCreatePasswordTokenUseCase(
	log contracts_providers.ILoggerProvider,
	passRepo contracts_repositories.IPasswordRepository,
	hashProvider contracts_providers.IHashProvider,
	repo contracts_repositories.IOneTimeTokenRepository,
) *CreatePasswordTokenUseCase {
	return &CreatePasswordTokenUseCase{
		BaseUseCaseValidation: usecase.BaseUseCaseValidation[dtos.PasswordTokenCreate, dtos.PasswordCreateNoHash]{
			AppMessages: locales.NewLocale(locales.EN_US),
			Guards:      usecase.NewGuards(),
		},
		log:              log,
		passRepo:         passRepo,
		hashProvider:     hashProvider,
		oneTimetokenRepo: repo,
	}
}
