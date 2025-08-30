package usecases_user

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

type ActivateUserUseCase struct {
	usecase.BaseUseCaseValidation[dtos.UserActivate, bool]
	log              contracts_providers.ILoggerProvider
	userRepo         contracts_repositories.IUserRepository
	oneTimetokenRepo contracts_repositories.IOneTimeTokenRepository

	hashProvider contracts_providers.IHashProvider
}

var _ usecase.BaseUseCase[dtos.UserActivate, bool] = (*ActivateUserUseCase)(nil)

func (uc *ActivateUserUseCase) SetLocale(locale locales.LocaleTypeEnum) {
	if locale != "" {
		uc.Locale = locale
	}
}

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
	_status := "active"
	updateUser.Status = &_status

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

func NewActivateUserUseCase(
	log contracts_providers.ILoggerProvider,
	userRepo contracts_repositories.IUserRepository,
	oneTimeTokenRepository contracts_repositories.IOneTimeTokenRepository,
	hashProvider contracts_providers.IHashProvider,
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
