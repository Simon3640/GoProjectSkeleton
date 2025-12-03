package integrationtest

import (
	"goprojectskeleton/src/application/shared/settings"
	"goprojectskeleton/src/infrastructure/providers"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEmailProvider_Integration(t *testing.T) {
	assert := assert.New(t)

	smtpHost := settings.AppSettingsInstance.MailHost
	smtpPort := settings.AppSettingsInstance.MailPort
	from := settings.AppSettingsInstance.MailFrom
	to := "recipient@example.com"
	subject := "Test Email"
	body := "This is a test email."

	email_provider := providers.NewEmailProvider()

	email_provider.Setup(smtpHost, smtpPort, from, "password")

	err := email_provider.SendEmail(to, subject, body)
	assert.Nil(err)
}
