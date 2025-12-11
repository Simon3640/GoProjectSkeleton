package userusecases

import (
	"context"

	contractsProviders "github.com/simon3640/goprojectskeleton/src/application/contracts/providers"
	usercontracts "github.com/simon3640/goprojectskeleton/src/application/modules/user/contracts"
	userdtos "github.com/simon3640/goprojectskeleton/src/application/modules/user/dtos"
	"github.com/simon3640/goprojectskeleton/src/application/shared/guards"
	"github.com/simon3640/goprojectskeleton/src/application/shared/locales"
	"github.com/simon3640/goprojectskeleton/src/application/shared/locales/messages"
	"github.com/simon3640/goprojectskeleton/src/application/shared/status"
	usecase "github.com/simon3640/goprojectskeleton/src/application/shared/use_case"
	"github.com/simon3640/goprojectskeleton/src/domain/models"
)

// UpdateUserUseCase is a use case that updates a user
type UpdateUserUseCase struct {
	usecase.BaseUseCaseValidation[userdtos.UserUpdate, models.User]
	log  contractsProviders.ILoggerProvider
	repo usercontracts.IUserRepository
}

var _ usecase.BaseUseCase[userdtos.UserUpdate, models.User] = (*UpdateUserUseCase)(nil)

// SetLocale sets the locale for the use case
func (uc *UpdateUserUseCase) SetLocale(locale locales.LocaleTypeEnum) {
	if locale != "" {
		uc.Locale = locale
	}
}

// Execute executes the use case
func (uc *UpdateUserUseCase) Execute(ctx context.Context,
	locale locales.LocaleTypeEnum,
	input userdtos.UserUpdate,
) *usecase.UseCaseResult[models.User] {
	result := usecase.NewUseCaseResult[models.User]()
	uc.SetLocale(locale)
	uc.Validate(ctx, input, result)
	if result.HasError() {
		return result
	}

	uc.updateUser(input, result)
	return result
}

// updateUser attempts to update the user.
// It sets errors in the result if the update fails; success data is set in the Execute method.
func (uc *UpdateUserUseCase) updateUser(input userdtos.UserUpdate, result *usecase.UseCaseResult[models.User]) {
	res, err := uc.repo.Update(input.ID, input)
	if err != nil {
		uc.log.Error("Error updating user", err.ToError())
		result.SetError(err.Code, uc.AppMessages.Get(uc.Locale, err.Context))
	}
	result.SetData(
		status.Success,
		*res,
		uc.AppMessages.Get(uc.Locale, messages.MessageKeysInstance.USER_WAS_CREATED))
}

// NewUpdateUserUseCase creates a new update user use case
func NewUpdateUserUseCase(
	log contractsProviders.ILoggerProvider,
	repo usercontracts.IUserRepository,
) *UpdateUserUseCase {
	return &UpdateUserUseCase{
		BaseUseCaseValidation: usecase.BaseUseCaseValidation[userdtos.UserUpdate, models.User]{
			AppMessages: locales.NewLocale(locales.EN_US),
			Guards:      usecase.NewGuards(guards.RoleGuard("admin", "user"), guards.UserResourceGuard[userdtos.UserUpdate]()),
		},
		log:  log,
		repo: repo,
	}
}
