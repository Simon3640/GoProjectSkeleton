package auth

import (
	"context"

	contracts_providers "gormgoskeleton/src/application/contracts/providers"
	contracts_repositories "gormgoskeleton/src/application/contracts/repositories"
	dtos "gormgoskeleton/src/application/shared/DTOs"
	"gormgoskeleton/src/application/shared/locales"
	"gormgoskeleton/src/application/shared/status"
	usecase "gormgoskeleton/src/application/shared/use_case"
	"gormgoskeleton/src/domain/models"
)

type GetResetPasswordTokenUseCase struct {
	usecase.BaseUseCaseValidation[string, string]
	log contracts_providers.ILoggerProvider

	tokenRepo contracts_repositories.IOneTimeTokenRepository
	userRepo  contracts_repositories.IUserRepository

	hashProvider contracts_providers.IHashProvider
}

func (uc *GetResetPasswordTokenUseCase) SetLocale(locale locales.LocaleTypeEnum) {
	if locale != "" {
		uc.Locale = locale
	}
}

// Execute generates a reset password token for the user identified by email or phone
func (uc *GetResetPasswordTokenUseCase) Execute(ctx context.Context,
	locale locales.LocaleTypeEnum,
	input string,
) *usecase.UseCaseResult[string] {
	result := usecase.NewUseCaseResult[string]()
	uc.SetLocale(locale)
	uc.Validate(ctx, input, result)
	if result.HasError() {
		return result
	}
	user, err := uc.userRepo.GetByEmailOrPhone(input)
	if err != nil {
		uc.log.Error("Error getting user by email or phone", err.ToError())
		result.SetError(
			err.Code,
			uc.AppMessages.Get(
				uc.Locale,
				err.Context,
			),
		)
	}

	token, hash, err := uc.hashProvider.OneTimeToken()
	if err != nil {
		uc.log.Error("Error generating one time token", err.ToError())
		result.SetError(
			err.Code,
			uc.AppMessages.Get(
				uc.Locale,
				err.Context,
			),
		)
	}

	tokenCreate := dtos.NewOneTimeTokenCreate(user.ID, models.OneTimeTokenPurposePasswordReset, hash)
	_, err = uc.tokenRepo.Create(*tokenCreate)
	if err != nil {
		uc.log.Error("Error creating one time token", err.ToError())
		result.SetError(
			err.Code,
			uc.AppMessages.Get(
				uc.Locale,
				err.Context,
			),
		)
	}

	result.SetData(
		status.Success,
		token,
		uc.AppMessages.Get(
			uc.Locale,
			"auth.reset_password_token_sent",
		),
	)

	return result
}

func NewGetResetPasswordTokenUseCase(
	log contracts_providers.ILoggerProvider,
	tokenRepo contracts_repositories.IOneTimeTokenRepository,
	userRepo contracts_repositories.IUserRepository,
	hashProvider contracts_providers.IHashProvider,
) *GetResetPasswordTokenUseCase {
	return &GetResetPasswordTokenUseCase{
		BaseUseCaseValidation: usecase.BaseUseCaseValidation[string, string]{
			AppMessages: locales.NewLocale(locales.EN_US),
			Guards:      usecase.NewGuards(),
		},
		log:          log,
		tokenRepo:    tokenRepo,
		userRepo:     userRepo,
		hashProvider: hashProvider,
	}
}
