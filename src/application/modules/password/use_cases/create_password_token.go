package usecases_password

import (
	"context"
	"time"

	contractsProviders "github.com/simon3640/goprojectskeleton/src/application/contracts/providers"
	contracts_repositories "github.com/simon3640/goprojectskeleton/src/application/contracts/repositories"
	dtos "github.com/simon3640/goprojectskeleton/src/application/shared/DTOs"
	"github.com/simon3640/goprojectskeleton/src/application/shared/locales"
	"github.com/simon3640/goprojectskeleton/src/application/shared/locales/messages"
	"github.com/simon3640/goprojectskeleton/src/application/shared/services"
	"github.com/simon3640/goprojectskeleton/src/application/shared/status"
	usecase "github.com/simon3640/goprojectskeleton/src/application/shared/use_case"
)

type CreatePasswordTokenUseCase struct {
	usecase.BaseUseCaseValidation[dtos.PasswordTokenCreate, bool]
	log              contractsProviders.ILoggerProvider
	passRepo         contracts_repositories.IPasswordRepository
	hashProvider     contractsProviders.IHashProvider
	oneTimetokenRepo contracts_repositories.IOneTimeTokenRepository
}

var _ usecase.BaseUseCase[dtos.PasswordTokenCreate, bool] = (*CreatePasswordTokenUseCase)(nil)

func (uc *CreatePasswordTokenUseCase) SetLocale(locale locales.LocaleTypeEnum) {
	if locale != "" {
		uc.Locale = locale
	}
}

func (uc *CreatePasswordTokenUseCase) Execute(ctx context.Context,
	locale locales.LocaleTypeEnum,
	input dtos.PasswordTokenCreate,
) *usecase.UseCaseResult[bool] {
	result := usecase.NewUseCaseResult[bool]()
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

	_, err = services.CreatePasswordService(passwordCreateNoHash, uc.hashProvider, uc.passRepo)

	if err != nil {
		uc.log.Error("CreatePasswordTokenUseCase: Execute: Error creating password", err.ToError())
		result.SetError(
			err.Code,
			uc.AppMessages.Get(
				uc.Locale,
				err.Context,
			),
		)
		return result
	}

	// Mark the token as used

	_, err = uc.oneTimetokenRepo.Update(oneTimeToke.ID,
		dtos.OneTimeTokenUpdate{IsUsed: true, ID: oneTimeToke.ID})

	if err != nil {
		uc.log.Error("Error updating one time token as used", err.ToError())
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
			messages.MessageKeysInstance.RESET_PASSWORD_TOKEN_VALID,
		),
	)

	return result
}

func NewCreatePasswordTokenUseCase(
	log contractsProviders.ILoggerProvider,
	passRepo contracts_repositories.IPasswordRepository,
	hashProvider contractsProviders.IHashProvider,
	repo contracts_repositories.IOneTimeTokenRepository,
) *CreatePasswordTokenUseCase {
	return &CreatePasswordTokenUseCase{
		BaseUseCaseValidation: usecase.BaseUseCaseValidation[dtos.PasswordTokenCreate, bool]{
			AppMessages: locales.NewLocale(locales.EN_US),
			Guards:      usecase.NewGuards(),
		},
		log:              log,
		passRepo:         passRepo,
		hashProvider:     hashProvider,
		oneTimetokenRepo: repo,
	}
}
