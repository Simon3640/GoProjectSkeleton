package templates

import "gormgoskeleton/src/application/shared/locales"

type TemplateKeysEnum string

type TemplateKeys struct {
	WelcomeEmail       TemplateKeysEnum
	PasswordResetEmail TemplateKeysEnum
}

var TemplateKeysInstance = TemplateKeys{
	WelcomeEmail:       "WELCOME_EMAIL",
	PasswordResetEmail: "PASSWORD_RESET_EMAIL",
}

var EnTemplates = map[TemplateKeysEnum]string{
	TemplateKeysInstance.WelcomeEmail:       "new_user_en.gohtml",
	TemplateKeysInstance.PasswordResetEmail: "password_reset_email.gohtml",
}

var EsTemplates = map[TemplateKeysEnum]string{
	TemplateKeysInstance.WelcomeEmail:       "new_user_es.gohtml",
	TemplateKeysInstance.PasswordResetEmail: "password_reset_email.gohtml",
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
