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
	"github.com/simon3640/goprojectskeleton/src/domain/models"
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

// Execute creates a password token for the user
// it creates a password token for the user and sets the password token in the result
// returns the password token if created successfully, otherwise returns an error
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
	oneTimeToken := uc.validateAndGetToken(result, input.Token)
	if result.HasError() {
		return result
	}

	uc.createPassword(result, oneTimeToken, input.NoHashedPassword)
	if result.HasError() {
		return result
	}

	uc.markTokenAsUsed(result, oneTimeToken.ID)
	if result.HasError() {
		return result
	}

	uc.setSuccessResult(result)
	return result
}

func (uc *CreatePasswordTokenUseCase) validateAndGetToken(result *usecase.UseCaseResult[bool], token string) *models.OneTimeToken {
	hash := uc.hashProvider.HashOneTimeToken(token)
	oneTimeToken, err := uc.oneTimetokenRepo.GetByTokenHash(hash)
	if err != nil {
		uc.log.Error("Error getting one time token by hash", err.ToError())
		result.SetError(
			err.Code,
			uc.AppMessages.Get(
				uc.Locale,
				err.Context,
			),
		)
		return nil
	}

	if oneTimeToken == nil || oneTimeToken.IsUsed || oneTimeToken.Expires.Before(time.Now()) {
		uc.log.Error("One time token is not valid", nil)
		result.SetError(
			status.Conflict,
			uc.AppMessages.Get(
				uc.Locale,
				messages.MessageKeysInstance.INVALID_PASSWORD_RESET_TOKEN,
			),
		)
		return nil
	}

	return oneTimeToken
}

func (uc *CreatePasswordTokenUseCase) createPassword(result *usecase.UseCaseResult[bool], oneTimeToken *models.OneTimeToken, noHashedPassword string) {
	passwordCreateNoHash := dtos.PasswordCreateNoHash{
		UserID:           oneTimeToken.UserID,
		NoHashedPassword: noHashedPassword,
		IsActive:         true,
	}

	_, err := services.CreatePasswordService(passwordCreateNoHash, uc.hashProvider, uc.passRepo)
	if err != nil {
		uc.log.Error("CreatePasswordTokenUseCase: Execute: Error creating password", err.ToError())
		result.SetError(
			err.Code,
			uc.AppMessages.Get(
				uc.Locale,
				err.Context,
			),
		)
		return
	}
}

func (uc *CreatePasswordTokenUseCase) markTokenAsUsed(result *usecase.UseCaseResult[bool], tokenID uint) {
	_, err := uc.oneTimetokenRepo.Update(tokenID,
		dtos.OneTimeTokenUpdate{IsUsed: true, ID: tokenID})

	if err != nil {
		uc.log.Error("Error updating one time token as used", err.ToError())
		result.SetError(
			err.Code,
			uc.AppMessages.Get(
				uc.Locale,
				err.Context,
			),
		)
		return
	}
}

func (uc *CreatePasswordTokenUseCase) setSuccessResult(result *usecase.UseCaseResult[bool]) {
	result.SetData(
		status.Success,
		true,
		uc.AppMessages.Get(
			uc.Locale,
			messages.MessageKeysInstance.RESET_PASSWORD_TOKEN_VALID,
		),
	)
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
