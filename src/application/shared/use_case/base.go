// Package usecase provides the base interface and validation for use cases.
package usecase

import (
	"strings"

	app_context "github.com/simon3640/goprojectskeleton/src/application/shared/context"
	"github.com/simon3640/goprojectskeleton/src/application/shared/locales"
	"github.com/simon3640/goprojectskeleton/src/application/shared/status"
)

// BaseUseCase is the base interface for all use cases
type BaseUseCase[Input any, Output any] interface {
	SetLocale(locale locales.LocaleTypeEnum)
	Execute(appContext *app_context.AppContext,
		locale locales.LocaleTypeEnum,
		input Input,
	) *UseCaseResult[Output]
	SetAppContext(appContext *app_context.AppContext)
}

// BaseUseCaseValidation is the base struct for all use cases with validation
// Abstracts the validation logic from the use case implementation
// Abstracts the guards logic from the use case implementation
type BaseUseCaseValidation[Input any, Output any] struct {
	Guards      Guards
	AppMessages *locales.Locale
	Locale      locales.LocaleTypeEnum
	AppContext  *app_context.AppContext
}

// SetAppContext sets the app context for the use case
func (v *BaseUseCaseValidation[Input, Output]) SetAppContext(appContext *app_context.AppContext) {
	if appContext != nil {
		v.AppContext = appContext
	} else {
		v.AppContext = app_context.NewVoidAppContext()
	}
}

// SetLocale sets the locale for the use case
func (v *BaseUseCaseValidation[Input, Output]) SetLocale(locale locales.LocaleTypeEnum) {
	if locale != "" {
		v.Locale = locale
	} else {
		v.Locale = locales.EN_US
	}
}

// Validate validates the input for the use case
// - If the input has the method Validate then call it
// - If the guards are empty then return
// - If the guards are not empty then set the actor and validate the input
// - If the input is invalid then set the error and return
// - If the input is valid then return
func (v *BaseUseCaseValidation[Input, Output]) Validate(
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
	v.Guards.SetActor(*v.AppContext.User)
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
