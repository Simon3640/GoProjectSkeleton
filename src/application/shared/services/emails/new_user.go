package email_service

import (
	email_models "goprojectskeleton/src/application/shared/services/emails/models"
)

type RegisterUserEmailService struct {
	EmailServiceBase[email_models.NewUserEmailData]
}

var RegisterUserEmailServiceInstance *RegisterUserEmailService

func init() {
	RegisterUserEmailServiceInstance = &RegisterUserEmailService{}
}
