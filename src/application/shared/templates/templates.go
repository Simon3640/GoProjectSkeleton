package templates

import "github.com/simon3640/goprojectskeleton/src/application/shared/locales"

type TemplateKeysEnum string

type TemplateKeys struct {
	WelcomeEmail       TemplateKeysEnum
	PasswordResetEmail TemplateKeysEnum
	OTPEmail           TemplateKeysEnum
}

var TemplateKeysInstance = TemplateKeys{
	WelcomeEmail:       "WELCOME_EMAIL",
	PasswordResetEmail: "PASSWORD_RESET_EMAIL",
	OTPEmail:           "OTP_EMAIL",
}

var EnTemplates = map[TemplateKeysEnum]string{
	TemplateKeysInstance.WelcomeEmail:       "new_user_en.gohtml",
	TemplateKeysInstance.PasswordResetEmail: "reset_password_en.gohtml",
	TemplateKeysInstance.OTPEmail:           "otp_en.gohtml",
}

var EsTemplates = map[TemplateKeysEnum]string{
	TemplateKeysInstance.WelcomeEmail:       "new_user_es.gohtml",
	TemplateKeysInstance.PasswordResetEmail: "reset_password_es.gohtml",
	TemplateKeysInstance.OTPEmail:           "otp_es.gohtml",
}

type Templates struct {
	En map[TemplateKeysEnum]string
	Es map[TemplateKeysEnum]string
}

var EmailTemplates = Templates{
	En: EnTemplates,
	Es: EsTemplates,
}

func GetTemplate(locale locales.LocaleTypeEnum, key TemplateKeysEnum) string {
	switch locale {
	case locales.EN_US:
		return EmailTemplates.En[key]
	case locales.ES_ES:
		return EmailTemplates.Es[key]
	default:
		return EmailTemplates.En[key]
	}
}
