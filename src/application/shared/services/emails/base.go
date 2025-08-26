package email_service

import (
	"gormgoskeleton/src/application/contracts"
	application_errors "gormgoskeleton/src/application/shared/errors"
	"gormgoskeleton/src/domain/models"
)

type EmailServiceBase[D any] struct {
	Renderer contracts.IRenderProvider[D]
	Sender   contracts.IEmailProvider
	template string
	subject  string
}

func (svc *EmailServiceBase[D]) SendWithTemplate(
	data D,
	user models.User,
) *application_errors.ApplicationError {
	rendered, err := svc.Renderer.Render(svc.template, data)
	if err != nil {
		return err
	}
	err = svc.Sender.SendEmail(user.Email, svc.subject, rendered)
	if err != nil {
		return err
	}
	return nil
}

const TemplatesPath = "src/infrastructure/templates/emails/"
