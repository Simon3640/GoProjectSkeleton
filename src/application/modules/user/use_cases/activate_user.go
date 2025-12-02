package userusecases

import (
	"context"
	"time"

	contractsProviders "goprojectskeleton/src/application/contracts/providers"
	contracts_repositories "goprojectskeleton/src/application/contracts/repositories"
	dtos "goprojectskeleton/src/application/shared/DTOs"
	"goprojectskeleton/src/application/shared/locales"
	"goprojectskeleton/src/application/shared/locales/messages"
	"goprojectskeleton/src/application/shared/status"
	usecase "goprojectskeleton/src/application/shared/use_case"
	"goprojectskeleton/src/domain/models"
)

// ActivateUserUseCase is a use case that activates a user
type ActivateUserUseCase struct {
	usecase.BaseUseCaseValidation[dtos.UserActivate, bool]
	log              contractsProviders.ILoggerProvider
	userRepo         contracts_repositories.IUserRepository
	oneTimetokenRepo contracts_repositories.IOneTimeTokenRepository

	hashProvider contractsProviders.IHashProvider
}

var _ usecase.BaseUseCase[dtos.UserActivate, bool] = (*ActivateUserUseCase)(nil)

// SetLocale sets the locale for the use case
func (uc *ActivateUserUseCase) SetLocale(locale locales.LocaleTypeEnum) {
	if locale != "" {
		uc.Locale = locale
	}
}

// Execute executes the use case
func (uc *ActivateUserUseCase) Execute(ctx context.Context,
	locale locales.LocaleTypeEnum,
	input dtos.UserActivate,
) *usecase.UseCaseResult[bool] {
	result := usecase.NewUseCaseResult[bool]()
	uc.SetLocale(locale)
	uc.Validate(ctx, input, result)
	if result.HasError() {
		return result
	}

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
		return result
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
		return result
	}

	updateUser := dtos.UserUpdate{}
	updateUser.ID = oneTimeToken.UserID
	userStatus := models.UserStatusActive
	updateUser.Status = &userStatus

	_, err = uc.userRepo.Update(oneTimeToken.UserID, updateUser)

	if err != nil {
		uc.log.Error("Error updating user", err.ToError())
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
			messages.MessageKeysInstance.USER_WAS_CREATED,
		),
	)
	return result
}

// NewActivateUserUseCase creates a new activate user use case
func NewActivateUserUseCase(
	log contractsProviders.ILoggerProvider,
	userRepo contracts_repositories.IUserRepository,
	oneTimeTokenRepository contracts_repositories.IOneTimeTokenRepository,
	hashProvider contractsProviders.IHashProvider,
) *ActivateUserUseCase {
	return &ActivateUserUseCase{
		BaseUseCaseValidation: usecase.BaseUseCaseValidation[dtos.UserActivate, bool]{
			AppMessages: locales.NewLocale(locales.EN_US),
			Guards:      usecase.NewGuards(),
		},
		log:              log,
		userRepo:         userRepo,
		oneTimetokenRepo: oneTimeTokenRepository,
		hashProvider:     hashProvider,
	}
}
