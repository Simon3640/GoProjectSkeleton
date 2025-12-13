// Package passwordusecases contains the use cases for the password module
package passwordusecases

import (
	"time"

	contractsproviders "github.com/simon3640/goprojectskeleton/src/application/contracts/providers"
	contractsrepositories "github.com/simon3640/goprojectskeleton/src/application/contracts/repositories"
	passwordcontracts "github.com/simon3640/goprojectskeleton/src/application/modules/password/contracts"
	dtos "github.com/simon3640/goprojectskeleton/src/application/modules/password/dtos"
	passwordservices "github.com/simon3640/goprojectskeleton/src/application/modules/password/services"
	shareddtos "github.com/simon3640/goprojectskeleton/src/application/shared/DTOs"
	app_context "github.com/simon3640/goprojectskeleton/src/application/shared/context"
	"github.com/simon3640/goprojectskeleton/src/application/shared/locales"
	"github.com/simon3640/goprojectskeleton/src/application/shared/locales/messages"
	"github.com/simon3640/goprojectskeleton/src/application/shared/status"
	usecase "github.com/simon3640/goprojectskeleton/src/application/shared/use_case"
	"github.com/simon3640/goprojectskeleton/src/domain/models"
)

// CreatePasswordTokenUseCase is the use case for creating a password token
type CreatePasswordTokenUseCase struct {
	usecase.BaseUseCaseValidation[dtos.PasswordTokenCreate, bool]
	log              contractsproviders.ILoggerProvider
	passRepo         passwordcontracts.IPasswordRepository
	hashProvider     contractsproviders.IHashProvider
	oneTimetokenRepo contractsrepositories.IOneTimeTokenRepository
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
func (uc *CreatePasswordTokenUseCase) Execute(ctx *app_context.AppContext,
	locale locales.LocaleTypeEnum,
	input dtos.PasswordTokenCreate,
) *usecase.UseCaseResult[bool] {
	result := usecase.NewUseCaseResult[bool]()
	uc.SetLocale(locale)
	uc.SetAppContext(ctx)
	uc.Validate(input, result)
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

	if oneTimeToken == nil || oneTimeToken.IsUsed || oneTimeToken.Expires.Before(time.Now()) ||
		oneTimeToken.Purpose != models.OneTimeTokenPurposePasswordReset {
		uc.log.Error("One time token is not valid or has incorrect purpose", nil)
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

	_, err := passwordservices.CreatePasswordService(passwordCreateNoHash, uc.hashProvider, uc.passRepo)
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
		shareddtos.OneTimeTokenUpdate{IsUsed: true, ID: tokenID})

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
	log contractsproviders.ILoggerProvider,
	passRepo passwordcontracts.IPasswordRepository,
	hashProvider contractsproviders.IHashProvider,
	repo contractsrepositories.IOneTimeTokenRepository,
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
