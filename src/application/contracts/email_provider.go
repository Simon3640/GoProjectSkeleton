package contracts

import application_errors "gormgoskeleton/src/application/shared/errors"

type IEmailProvider interface {
	SendEmail(to string, subject string, body string) *application_errors.ApplicationError
}
