package userusecases

import (
	usercontracts "github.com/simon3640/goprojectskeleton/src/application/modules/user/contracts"
	userdtos "github.com/simon3640/goprojectskeleton/src/application/modules/user/dtos"
	app_context "github.com/simon3640/goprojectskeleton/src/application/shared/context"
	"github.com/simon3640/goprojectskeleton/src/application/shared/guards"
	"github.com/simon3640/goprojectskeleton/src/application/shared/locales"
	"github.com/simon3640/goprojectskeleton/src/application/shared/locales/messages"
	"github.com/simon3640/goprojectskeleton/src/application/shared/observability"
	"github.com/simon3640/goprojectskeleton/src/application/shared/status"
	usecase "github.com/simon3640/goprojectskeleton/src/application/shared/use_case"
	usermodels "github.com/simon3640/goprojectskeleton/src/domain/user/models"
)

// UpdateUserUseCase is a use case that updates a user
type UpdateUserUseCase struct {
	usecase.BaseUseCaseValidation[userdtos.UserUpdate, usermodels.User]
	repo usercontracts.IUserRepository
}

var _ usecase.BaseUseCase[userdtos.UserUpdate, usermodels.User] = (*UpdateUserUseCase)(nil)

// Execute executes the use case
func (uc *UpdateUserUseCase) Execute(ctx *app_context.AppContext,
	locale locales.LocaleTypeEnum,
	input userdtos.UserUpdate,
) *usecase.UseCaseResult[usermodels.User] {
	result := usecase.NewUseCaseResult[usermodels.User]()
	uc.SetLocale(locale)
	uc.SetAppContext(ctx)
	uc.Validate(input, result)
	if result.HasError() {
		return result
	}

	uc.updateUser(input, result)
	if result.HasError() {
		return result
	}
	observability.GetObservabilityComponents().Logger.InfoWithContext("User updated successfully", uc.AppContext)
	return result
}

// updateUser attempts to update the user.
// It sets errors in the result if the update fails; success data is set in the Execute method.
func (uc *UpdateUserUseCase) updateUser(input userdtos.UserUpdate, result *usecase.UseCaseResult[usermodels.User]) {
	res, err := uc.repo.Update(input.ID, input)
	if err != nil {
		observability.GetObservabilityComponents().Logger.ErrorWithContext("Error updating user", err.ToError(), uc.AppContext)
		result.SetError(err.Code, uc.AppMessages.Get(uc.Locale, err.Context))
	}
	result.SetData(
		status.Updated,
		*res,
		uc.AppMessages.Get(uc.Locale, messages.MessageKeysInstance.USER_WAS_CREATED))
}

// NewUpdateUserUseCase creates a new update user use case
func NewUpdateUserUseCase(
	repo usercontracts.IUserRepository,
) *UpdateUserUseCase {
	return &UpdateUserUseCase{
		BaseUseCaseValidation: usecase.BaseUseCaseValidation[userdtos.UserUpdate, usermodels.User]{
			AppMessages: locales.NewLocale(locales.EN_US),
			Guards:      usecase.NewGuards(guards.RoleGuard("admin", "user"), guards.UserResourceGuard[userdtos.UserUpdate]()),
		},
		repo: repo,
	}
}
