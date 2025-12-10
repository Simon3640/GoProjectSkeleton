package authusecases

import (
	"context"

	contractsProviders "github.com/simon3640/goprojectskeleton/src/application/contracts/providers"
	contracts_repositories "github.com/simon3640/goprojectskeleton/src/application/contracts/repositories"
	dtos "github.com/simon3640/goprojectskeleton/src/application/shared/DTOs"
	"github.com/simon3640/goprojectskeleton/src/application/shared/locales"
	"github.com/simon3640/goprojectskeleton/src/application/shared/locales/messages"
	"github.com/simon3640/goprojectskeleton/src/application/shared/status"
	usecase "github.com/simon3640/goprojectskeleton/src/application/shared/use_case"
	"github.com/simon3640/goprojectskeleton/src/domain/models"
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

	user := uc.getUser(result, input)
	if result.HasError() {
		return result
	}

	token, hash := uc.generateToken(result)
	if result.HasError() {
		return result
	}

	uc.createToken(result, user.ID, hash)
	if result.HasError() {
		return result
	}

	uc.setSuccessResult(result, user, token)
	return result
}

func (uc *GetResetPasswordTokenUseCase) getUser(result *usecase.UseCaseResult[dtos.OneTimeTokenUser], emailOrPhone string) *models.User {
	user, err := uc.userRepo.GetByEmailOrPhone(emailOrPhone)
	if err != nil {
		uc.log.Error("Error getting user by email or phone", err.ToError())
		result.SetError(
			err.Code,
			uc.AppMessages.Get(
				uc.Locale,
				err.Context,
			),
		)
		return nil
	}
	return user
}

func (uc *GetResetPasswordTokenUseCase) generateToken(result *usecase.UseCaseResult[dtos.OneTimeTokenUser]) (string, []byte) {
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
		return "", nil
	}
	return token, hash
}

func (uc *GetResetPasswordTokenUseCase) createToken(result *usecase.UseCaseResult[dtos.OneTimeTokenUser], userID uint, hash []byte) {
	tokenCreate := dtos.NewOneTimeTokenCreate(userID, models.OneTimeTokenPurposePasswordReset, hash)
	_, err := uc.tokenRepo.Create(*tokenCreate)
	if err != nil {
		uc.log.Error("Error creating one time token", err.ToError())
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

func (uc *GetResetPasswordTokenUseCase) setSuccessResult(result *usecase.UseCaseResult[dtos.OneTimeTokenUser], user *models.User, token string) {
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
