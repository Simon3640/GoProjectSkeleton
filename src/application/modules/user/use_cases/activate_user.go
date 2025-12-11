package userusecases

import (
	"context"
	"time"

	contractsProviders "github.com/simon3640/goprojectskeleton/src/application/contracts/providers"
	contractrepositories "github.com/simon3640/goprojectskeleton/src/application/contracts/repositories"
	usercontracts "github.com/simon3640/goprojectskeleton/src/application/modules/user/contracts"
	userdtos "github.com/simon3640/goprojectskeleton/src/application/modules/user/dtos"
	"github.com/simon3640/goprojectskeleton/src/application/shared/locales"
	"github.com/simon3640/goprojectskeleton/src/application/shared/locales/messages"
	"github.com/simon3640/goprojectskeleton/src/application/shared/status"
	usecase "github.com/simon3640/goprojectskeleton/src/application/shared/use_case"
	"github.com/simon3640/goprojectskeleton/src/domain/models"
)

// ActivateUserUseCase is a use case that activates a user
type ActivateUserUseCase struct {
	usecase.BaseUseCaseValidation[userdtos.UserActivate, bool]
	log              contractsProviders.ILoggerProvider
	userRepo         usercontracts.IUserRepository
	oneTimetokenRepo contractrepositories.IOneTimeTokenRepository

	hashProvider contractsProviders.IHashProvider
}

var _ usecase.BaseUseCase[userdtos.UserActivate, bool] = (*ActivateUserUseCase)(nil)

// SetLocale sets the locale for the use case
func (uc *ActivateUserUseCase) SetLocale(locale locales.LocaleTypeEnum) {
	if locale != "" {
		uc.Locale = locale
	}
}

// Execute executes the use case
func (uc *ActivateUserUseCase) Execute(ctx context.Context,
	locale locales.LocaleTypeEnum,
	input userdtos.UserActivate,
) *usecase.UseCaseResult[bool] {
	result := usecase.NewUseCaseResult[bool]()
	uc.SetLocale(locale)
	uc.Validate(ctx, input, result)
	if result.HasError() {
		return result
	}

	userID := uc.validateOneTimeToken(input, result)
	if result.HasError() {
		return result
	}

	uc.updateUser(*userID, result)
	if result.HasError() {
		return result
	}

	result.SetData(
		status.Success,
		true,
		uc.AppMessages.Get(
			uc.Locale,
			messages.MessageKeysInstance.USER_WAS_CREATED,
		),
	)
	return result
}

// validateOneTimeToken validates the one time token
// returns the user id if the token is valid
func (uc *ActivateUserUseCase) validateOneTimeToken(input userdtos.UserActivate, result *usecase.UseCaseResult[bool]) *uint {
	hash := uc.hashProvider.HashOneTimeToken(input.Token)
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
				messages.MessageKeysInstance.INVALID_USER_ACTIVATION_TOKEN,
			),
		)
		return nil
	}
	return &oneTimeToken.UserID
}

// updateUser updates the user status to active
func (uc *ActivateUserUseCase) updateUser(userID uint, result *usecase.UseCaseResult[bool]) {
	updateUser := userdtos.UserUpdate{}
	updateUser.ID = userID
	userStatus := models.UserStatusActive
	updateUser.Status = &userStatus

	_, err := uc.userRepo.Update(userID, updateUser)

	if err != nil {
		uc.log.Error("Error updating user", err.ToError())
		result.SetError(
			err.Code,
			uc.AppMessages.Get(
				uc.Locale,
				err.Context,
			),
		)
	}
}

// NewActivateUserUseCase creates a new activate user use case
func NewActivateUserUseCase(
	log contractsProviders.ILoggerProvider,
	userRepo usercontracts.IUserRepository,
	oneTimeTokenRepository contractrepositories.IOneTimeTokenRepository,
	hashProvider contractsProviders.IHashProvider,
) *ActivateUserUseCase {
	return &ActivateUserUseCase{
		BaseUseCaseValidation: usecase.BaseUseCaseValidation[userdtos.UserActivate, bool]{
			AppMessages: locales.NewLocale(locales.EN_US),
			Guards:      usecase.NewGuards(),
		},
		log:              log,
		userRepo:         userRepo,
		oneTimetokenRepo: oneTimeTokenRepository,
		hashProvider:     hashProvider,
	}
}
