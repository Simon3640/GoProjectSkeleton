package email_service

import (
	email_models "github.com/simon3640/goprojectskeleton/src/application/shared/services/emails/models"
)

type OneTimePasswordEmailService struct {
	EmailServiceBase[email_models.OneTimePasswordEmailData]
}

var OneTimePasswordEmailServiceInstance *OneTimePasswordEmailService

func init() {
	OneTimePasswordEmailServiceInstance = &OneTimePasswordEmailService{}
}
