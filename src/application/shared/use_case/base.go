// Package usecase provides the base interface and validation for use cases.
package usecase

import (
	"strings"

	app_context "github.com/simon3640/goprojectskeleton/src/application/shared/context"
	"github.com/simon3640/goprojectskeleton/src/application/shared/locales"
	"github.com/simon3640/goprojectskeleton/src/application/shared/status"
)

type BaseUseCase[Input any, Output any] interface {
	SetLocale(locale locales.LocaleTypeEnum)
	Execute(appContext *app_context.AppContext,
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
	appContext *app_context.AppContext,
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
	v.Guards.SetActor(*appContext.User)
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
