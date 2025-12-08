package email_service

import (
	email_models "github.com/simon3640/goprojectskeleton/src/application/shared/services/emails/models"
)

type ResetPasswordEmailService struct {
	EmailServiceBase[email_models.ResetPasswordEmailData]
}

var ResetPasswordEmailServiceInstance *ResetPasswordEmailService

func init() {
	ResetPasswordEmailServiceInstance = &ResetPasswordEmailService{}
}
