package auth

import (
	"context"

	contractsProviders "gormgoskeleton/src/application/contracts/providers"
	contracts_repositories "gormgoskeleton/src/application/contracts/repositories"
	dtos "gormgoskeleton/src/application/shared/DTOs"
	"gormgoskeleton/src/application/shared/locales"
	"gormgoskeleton/src/application/shared/locales/messages"
	"gormgoskeleton/src/application/shared/status"
	usecase "gormgoskeleton/src/application/shared/use_case"
	"gormgoskeleton/src/domain/models"
)

type GetResetPasswordTokenUseCase struct {
	usecase.BaseUseCaseValidation[string, dtos.OneTimeTokenUser]
	log contractsProviders.ILoggerProvider

	tokenRepo contracts_repositories.IOneTimeTokenRepository
	userRepo  contracts_repositories.IUserRepository

	hashProvider contractsProviders.IHashProvider
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
) *usecase.UseCaseResult[dtos.OneTimeTokenUser] {
	result := usecase.NewUseCaseResult[dtos.OneTimeTokenUser]()
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
		return result
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
		return result
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
		return result
	}

	result.SetData(
		status.Success,
		dtos.OneTimeTokenUser{
			User:  *user,
			Token: token,
		},
		uc.AppMessages.Get(
			uc.Locale,
			messages.MessageKeysInstance.PASSWORD_TOKEN_CREATED,
		),
	)

	return result
}

func (uc *GetResetPasswordTokenUseCase) Validate(
	ctx context.Context,
	input string,
	result *usecase.UseCaseResult[dtos.OneTimeTokenUser],
) {
	// Skip input validation as it's just a string (email or phone)
}

func NewGetResetPasswordTokenUseCase(
	log contractsProviders.ILoggerProvider,
	tokenRepo contracts_repositories.IOneTimeTokenRepository,
	userRepo contracts_repositories.IUserRepository,
	hashProvider contractsProviders.IHashProvider,
) *GetResetPasswordTokenUseCase {
	return &GetResetPasswordTokenUseCase{
		BaseUseCaseValidation: usecase.BaseUseCaseValidation[string, dtos.OneTimeTokenUser]{
			AppMessages: locales.NewLocale(locales.EN_US),
			Guards:      usecase.NewGuards(),
		},
		log:          log,
		tokenRepo:    tokenRepo,
		userRepo:     userRepo,
		hashProvider: hashProvider,
	}
}
