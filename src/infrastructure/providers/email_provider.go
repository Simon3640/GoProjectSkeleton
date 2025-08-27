package providers

import (
	"net/smtp"
	"strconv"

	contracts_providers "gormgoskeleton/src/application/contracts/providers"
	application_errors "gormgoskeleton/src/application/shared/errors"
	"gormgoskeleton/src/application/shared/locales/messages"
	"gormgoskeleton/src/application/shared/settings"
	"gormgoskeleton/src/application/shared/status"
)

type EmailProvider struct {
	smtpHost string
	smtpPort int
	from     string
	password string
}

var _ contracts_providers.IEmailProvider = (*EmailProvider)(nil)

func (ep *EmailProvider) Setup(smtpHost string, smtpPort int, from string, password string) {
	ep.smtpHost = smtpHost
	ep.smtpPort = smtpPort
	ep.from = from
	ep.password = password
}

func (ep *EmailProvider) buildConnection() func(to string, message string) error {
	return func(to string, message string) error {
		var auth smtp.Auth
		// for testing purposes
		if settings.AppSettingsInstance.AppEnv == "development" || settings.AppSettingsInstance.AppEnv == "test" {
			auth = nil
		} else {
			auth = smtp.PlainAuth("", ep.from, ep.password, ep.smtpHost)
		}
		err := smtp.SendMail(
			ep.smtpHost+":"+strconv.Itoa(ep.smtpPort),
			auth,
			ep.from,
			[]string{to},
			[]byte(message),
		)
		return err
	}
}

func (ep *EmailProvider) SendEmail(
	to string,
	subject string,
	body string,
) *application_errors.ApplicationError {
	message := "From: " + ep.from + "\n" +
		"To: " + to + "\n" +
		"Subject: " + subject + "\n" +
		"\n" +
		body
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
