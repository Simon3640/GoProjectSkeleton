package email_service

import (
	contracts_providers "gormgoskeleton/src/application/contracts/providers"
	email_models "gormgoskeleton/src/application/shared/services/emails/models"
)

type RegisterUserEmailService struct {
	EmailServiceBase[email_models.NewUserEmailData]
}

func (svc *RegisterUserEmailService) SetUp(
	renderer contracts_providers.IRendererProvider[email_models.NewUserEmailData],
	sender contracts_providers.IEmailProvider,
) {
	svc.EmailServiceBase.Renderer = renderer
	svc.EmailServiceBase.Sender = sender
}

var RegisterUserEmailServiceInstance *RegisterUserEmailService

func init() {
	RegisterUserEmailServiceInstance = &RegisterUserEmailService{}
}
