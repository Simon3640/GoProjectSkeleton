package contractsproviders

import application_errors "goprojectskeleton/src/application/shared/errors"

type IEmailProvider interface {
	SendEmail(to string, subject string, body string) *application_errors.ApplicationError
}
