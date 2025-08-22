package usecase

import (
	"context"
	"strings"

	app_context "gormgoskeleton/src/application/shared/context"
	"gormgoskeleton/src/application/shared/locales"
	"gormgoskeleton/src/application/shared/locales/messages"
	"gormgoskeleton/src/application/shared/status"
	"gormgoskeleton/src/domain/models"
)

type BaseUseCase[Input any, Output any] interface {
	SetLocale(locale locales.LocaleTypeEnum)
	Execute(ctx context.Context,
		locale locales.LocaleTypeEnum,
		input Input,
	) *UseCaseResult[Output]
}

type BaseUseCaseValidation[Input any, Output any] struct {
	Guards      Guards
	AppMessages *locales.Locale
	Locale      locales.LocaleTypeEnum
}

func (v *BaseUseCaseValidation[Input, Output]) Validate(
	ctx context.Context,
	input Input,
	result *UseCaseResult[Output],
) {
	// Know if input has the method Validation then call
	if validator, ok := any(input).(interface{ Validate() []string }); ok {
		errs := validator.Validate()
		if len(errs) > 0 {
			result.SetError(
				status.InvalidInput,
				strings.Join(errs, "\n"),
			)
			return
		}
	}

	user_ctx, ok := ctx.Value(app_context.UserKey).(models.UserWithRole)
	if !ok {
		result.SetError(
			status.Unauthorized,
			v.AppMessages.Get(
				v.Locale,
				messages.MessageKeysInstance.AUTHORIZATION_REQUIRED,
			),
		)
		return
	}
	v.Guards.SetActor(user_ctx)
	if err := v.Guards.Validate(input); err != nil {
		result.SetError(
			status.Unauthorized,
			v.AppMessages.Get(
				v.Locale,
				messages.MessageKeysInstance.UNAUTHORIZED_RESOURCE,
			),
		)
		return
	}
}
