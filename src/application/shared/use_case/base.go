package usecase

import (
	"context"
	"strings"

	app_context "github.com/simon3640/goprojectskeleton/src/application/shared/context"
	"github.com/simon3640/goprojectskeleton/src/application/shared/locales"
	"github.com/simon3640/goprojectskeleton/src/application/shared/locales/messages"
	"github.com/simon3640/goprojectskeleton/src/application/shared/status"
	"github.com/simon3640/goprojectskeleton/src/domain/models"
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
	if len(v.Guards.list) == 0 {
		return
	}
	user := ctx.Value(app_context.UserKey)
	if user == nil {
		result.SetError(
			status.InternalError,
			v.AppMessages.Get(
				v.Locale,
				messages.MessageKeysInstance.SOMETHING_WENT_WRONG,
			),
		)
		return
	}
	user_ctx, ok := user.(models.UserWithRole)
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
				*err,
			),
		)
		return
	}
}
