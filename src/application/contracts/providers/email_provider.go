package contractsproviders

import application_errors "github.com/simon3640/goprojectskeleton/src/application/shared/errors"

type IEmailProvider interface {
	SendEmail(to string, subject string, body string) *application_errors.ApplicationError
}
