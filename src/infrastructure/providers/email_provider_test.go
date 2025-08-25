package providers

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEmailProvider_Integration(t *testing.T) {
	assert := assert.New(t)

	smtpHost := "localhost"
	smtpPort := 1025
	from := "noreply@goprojectekeleton.com"
	to := "recipient@example.com"
	subject := "Test Email"
	body := "This is a test email."

	email_provider := NewEmailProvider(smtpHost, smtpPort, from, "password")

	err := email_provider.SendEmail(to, subject, body)
	assert.Nil(err)
}
