package providers

import (
	"fmt"
	"net/smtp"
	"strconv"

	contractsProviders "goprojectskeleton/src/application/contracts/providers"
	application_errors "goprojectskeleton/src/application/shared/errors"
	"goprojectskeleton/src/application/shared/locales/messages"
	"goprojectskeleton/src/application/shared/settings"
	"goprojectskeleton/src/application/shared/status"
)

type EmailProvider struct {
	smtpHost string
	smtpPort int
	from     string
	password string
}

var _ contractsProviders.IEmailProvider = (*EmailProvider)(nil)

func (ep *EmailProvider) Setup(smtpHost string, smtpPort int, from string, password string) {
	ep.smtpHost = smtpHost
	ep.smtpPort = smtpPort
	ep.from = from
	ep.password = password
}

func (ep *EmailProvider) buildConnection() func(to string, message []byte) error {
	return func(to string, message []byte) error {
		var auth smtp.Auth
		// para testing con Mailpit: no hace falta autenticación
		if !settings.AppSettingsInstance.MailAuthRequired {
			auth = nil
		} else {
			auth = smtp.PlainAuth("", ep.from, ep.password, ep.smtpHost)
		}
		err := smtp.SendMail(
			ep.smtpHost+":"+strconv.Itoa(ep.smtpPort),
			auth,
			ep.from,
			[]string{to},
			message,
		)
		return err
	}
}

func (ep *EmailProvider) SendEmail(
	to string,
	subject string,
	body string,
) *application_errors.ApplicationError {

	// Headers + cuerpo en HTML
	message := []byte(fmt.Sprintf(
		"From: %s\r\n"+
			"To: %s\r\n"+
			"Subject: %s\r\n"+
			"MIME-Version: 1.0\r\n"+
			"Content-Type: text/html; charset=\"UTF-8\"\r\n"+
			"\r\n"+
			"%s\r\n",
		ep.from, to, subject, body,
	))

	if err := ep.buildConnection()(to, message); err != nil {
		return application_errors.NewApplicationError(
			status.ProviderError,
			messages.MessageKeysInstance.SOMETHING_WENT_WRONG,
			err.Error(),
		)
	}
	return nil
}

func NewEmailProvider() *EmailProvider {
	return &EmailProvider{}
}

var EmailProviderInstance *EmailProvider

func init() {
	EmailProviderInstance = NewEmailProvider()
}
