package authusecases

import (
	contractproviders "github.com/simon3640/goprojectskeleton/src/application/contracts/providers"
	contractrepositories "github.com/simon3640/goprojectskeleton/src/application/contracts/repositories"
	authcontracts "github.com/simon3640/goprojectskeleton/src/application/modules/auth/contracts"
	dtos "github.com/simon3640/goprojectskeleton/src/application/shared/DTOs"
	app_context "github.com/simon3640/goprojectskeleton/src/application/shared/context"
	"github.com/simon3640/goprojectskeleton/src/application/shared/locales"
	"github.com/simon3640/goprojectskeleton/src/application/shared/locales/messages"
	"github.com/simon3640/goprojectskeleton/src/application/shared/observability"
	"github.com/simon3640/goprojectskeleton/src/application/shared/status"
	usecase "github.com/simon3640/goprojectskeleton/src/application/shared/use_case"
	"github.com/simon3640/goprojectskeleton/src/domain/models"
)

// GetResetPasswordTokenUseCase is the use case for generating a reset password token
type GetResetPasswordTokenUseCase struct {
	usecase.BaseUseCaseValidation[string, bool]

	tokenRepo contractrepositories.IOneTimeTokenRepository
	userRepo  authcontracts.IUserRepository

	hashProvider contractproviders.IHashProvider
}

var _ usecase.BaseUseCase[string, bool] = (*GetResetPasswordTokenUseCase)(nil)

// SetLocale sets the locale for the use case
func (uc *GetResetPasswordTokenUseCase) SetLocale(locale locales.LocaleTypeEnum) {
	if locale != "" {
		uc.Locale = locale
	}
}

// Execute generates a reset password token for the user identified by email or phone
func (uc *GetResetPasswordTokenUseCase) Execute(ctx *app_context.AppContext,
	locale locales.LocaleTypeEnum,
	input string,
) *usecase.UseCaseResult[bool] {
	result := usecase.NewUseCaseResult[bool]()
	uc.SetLocale(locale)
	uc.SetAppContext(ctx)
	uc.Validate(input, result)
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
	observability.GetObservabilityComponents().Logger.InfoWithContext("Reset password token created successfully", uc.AppContext)
	return result
}

func (uc *GetResetPasswordTokenUseCase) getUser(result *usecase.UseCaseResult[bool], emailOrPhone string) *models.User {
	user, err := uc.userRepo.GetByEmailOrPhone(emailOrPhone)
	if err != nil {
		observability.GetObservabilityComponents().Logger.ErrorWithContext("Error getting user by email or phone", err.ToError(), uc.AppContext)
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

func (uc *GetResetPasswordTokenUseCase) generateToken(result *usecase.UseCaseResult[bool]) (string, []byte) {
	token, hash, err := uc.hashProvider.OneTimeToken()
	if err != nil {
		observability.GetObservabilityComponents().Logger.ErrorWithContext("Error generating one time token", err.ToError(), uc.AppContext)
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

func (uc *GetResetPasswordTokenUseCase) createToken(result *usecase.UseCaseResult[bool], userID uint, hash []byte) {
	tokenCreate := dtos.NewOneTimeTokenCreate(userID, models.OneTimeTokenPurposePasswordReset, hash)
	_, err := uc.tokenRepo.Create(*tokenCreate)
	if err != nil {
		observability.GetObservabilityComponents().Logger.ErrorWithContext("Error creating one time token", err.ToError(), uc.AppContext)
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

func (uc *GetResetPasswordTokenUseCase) setSuccessResult(result *usecase.UseCaseResult[bool], user *models.User, token string) {
	oneTimeToken := dtos.OneTimeTokenUser{User: *user, Token: token}
	uc.AppContext.AddOneTimeTokenToContext(oneTimeToken)
	result.SetData(
		status.Success,
		true,
		uc.AppMessages.Get(
			uc.Locale,
			messages.MessageKeysInstance.PASSWORD_TOKEN_CREATED,
		),
	)
}

func NewGetResetPasswordTokenUseCase(
	tokenRepo contractrepositories.IOneTimeTokenRepository,
	userRepo authcontracts.IUserRepository,
	hashProvider contractproviders.IHashProvider,
) *GetResetPasswordTokenUseCase {
	return &GetResetPasswordTokenUseCase{
		BaseUseCaseValidation: usecase.BaseUseCaseValidation[string, bool]{
			AppMessages: locales.NewLocale(locales.EN_US),
			Guards:      usecase.NewGuards(),
		},
		tokenRepo:    tokenRepo,
		userRepo:     userRepo,
		hashProvider: hashProvider,
	}
}
