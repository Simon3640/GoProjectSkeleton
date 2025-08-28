package email_service

import (
	email_models "gormgoskeleton/src/application/shared/services/emails/models"
)

type ResetPasswordEmailService struct {
	EmailServiceBase[email_models.ResetPasswordEmailData]
}

var ResetPasswordEmailServiceInstance *ResetPasswordEmailService

func init() {
	ResetPasswordEmailServiceInstance = &ResetPasswordEmailService{}
}
