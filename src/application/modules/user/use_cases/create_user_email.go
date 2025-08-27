package usecases_user

import (
	"context"
	"fmt"

	"gormgoskeleton/src/application/contracts"
	"gormgoskeleton/src/application/shared/locales"
	"gormgoskeleton/src/application/shared/locales/messages"
	email_service "gormgoskeleton/src/application/shared/services/emails"
	email_models "gormgoskeleton/src/application/shared/services/emails/models"
	"gormgoskeleton/src/application/shared/settings"
	"gormgoskeleton/src/application/shared/status"
	usecase "gormgoskeleton/src/application/shared/use_case"
	"gormgoskeleton/src/domain/models"
)

type CreateUserSendEmailUseCase struct {
	appMessages *locales.Locale
	log         contracts.ILoggerProvider
	jwt         contracts.IJWTProvider
	locale      locales.LocaleTypeEnum
}

var _ usecase.BaseUseCase[models.User, models.User] = (*CreateUserSendEmailUseCase)(nil)

func (uc *CreateUserSendEmailUseCase) SetLocale(locale locales.LocaleTypeEnum) {
	if locale != "" {
		uc.locale = locale
	}
}

func (uc *CreateUserSendEmailUseCase) Execute(ctx context.Context,
	locale locales.LocaleTypeEnum,
	input models.User,
) *usecase.UseCaseResult[models.User] {
	result := usecase.NewUseCaseResult[models.User]()
	uc.SetLocale(locale)

	token, exp, err := uc.jwt.GenerateAccessToken(ctx, fmt.Sprint(input.ID), nil)

	newUserEmailData := email_models.NewUserEmailData{
		Name:            input.Name,
		ActivationToken: token,
		Expiration:      exp,
		AppName:         settings.AppSettingsInstance.AppName,
		SupportEmail:    settings.AppSettingsInstance.AppSupportEmail,
	}

	if err := email_service.RegisterUserEmailServiceInstance.SendWithTemplate(newUserEmailData, input, locale); err != nil {
		uc.log.Error("Error sending email", err.ToError())
		result.SetError(
			err.Code,
			uc.appMessages.Get(
				uc.locale,
				err.Context,
			),
		)
		return result
	}

	if err != nil {
		uc.log.Error("Error generating token", err.ToError())
		result.SetError(
			err.Code,
			uc.appMessages.Get(
				uc.locale,
				err.Context,
			),
		)
		return result
	}
	result.SetData(
		status.Success,
		input,
		uc.appMessages.Get(
			uc.locale,
			messages.MessageKeysInstance.USER_WAS_CREATED,
		),
	)
	return result
}

func NewCreateUserSendEmailUseCase(
	log contracts.ILoggerProvider,
	jwt contracts.IJWTProvider,
) *CreateUserSendEmailUseCase {
	return &CreateUserSendEmailUseCase{
		appMessages: locales.NewLocale(locales.EN_US),
		log:         log,
		jwt:         jwt,
	}
}
