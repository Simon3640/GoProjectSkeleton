package email_service

import (
	email_models "gormgoskeleton/src/application/shared/services/emails/models"
)

type OneTimePasswordEmailService struct {
	EmailServiceBase[email_models.OneTimePasswordEmailData]
}

var OneTimePasswordEmailServiceInstance *OneTimePasswordEmailService

func init() {
	OneTimePasswordEmailServiceInstance = &OneTimePasswordEmailService{}
}
