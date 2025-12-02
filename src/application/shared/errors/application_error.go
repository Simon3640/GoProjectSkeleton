package application_errors

import (
	"errors"
	"fmt"

	"goprojectskeleton/src/application/shared/locales/messages"
	"goprojectskeleton/src/application/shared/status"
)

type ApplicationError struct {
	Code    status.ApplicationStatusEnum
	Context messages.MessageKeysEnum
	ErrMsg  string
}

func (ae *ApplicationError) ToError() error {
	errorMessage := fmt.Sprintf(
		"Error: %s, Context: %s, Code: %s",
		ae.ErrMsg,
		ae.Context,
		string(ae.Code),
	)
	return errors.New(errorMessage)
}

func NewApplicationError(code status.ApplicationStatusEnum, context messages.MessageKeysEnum, errMsg string) *ApplicationError {
	return &ApplicationError{
		Code:    code,
		Context: context,
		ErrMsg:  errMsg,
	}
}
