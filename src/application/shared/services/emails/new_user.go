package email_service

import (
	"gormgoskeleton/src/application/contracts"
	email_models "gormgoskeleton/src/application/shared/services/emails/models"
)

type RegisterUserEmailService struct {
	EmailServiceBase[email_models.NewUserEmailData]
}

func (svc *RegisterUserEmailService) SetUp(
	renderer contracts.IRendererProvider[email_models.NewUserEmailData],
	sender contracts.IEmailProvider,
) {
	svc.EmailServiceBase.Renderer = renderer
	svc.EmailServiceBase.Sender = sender
}

var RegisterUserEmailServiceInstance *RegisterUserEmailService

func init() {
	RegisterUserEmailServiceInstance = &RegisterUserEmailService{}
}
