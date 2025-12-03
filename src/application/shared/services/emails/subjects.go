package email_service

import "goprojectskeleton/src/application/shared/locales"

type SubjectKeysEnum string

type SubjectKeys struct {
	WelcomeEmail       SubjectKeysEnum
	PasswordResetEmail SubjectKeysEnum
	OTPEmail           SubjectKeysEnum
}

var SubjectKeysInstance = SubjectKeys{
	WelcomeEmail:       "WELCOME_EMAIL",
	PasswordResetEmail: "PASSWORD_RESET_EMAIL",
	OTPEmail:           "OTP_EMAIL",
}

var EnSubjects = map[SubjectKeysEnum]string{
	SubjectKeysInstance.WelcomeEmail:       "Welcome to our App!",
	SubjectKeysInstance.PasswordResetEmail: "Password Reset Instructions",
	SubjectKeysInstance.OTPEmail:           "Your One-Time Password (OTP)",
}

var EsSubjects = map[SubjectKeysEnum]string{
	SubjectKeysInstance.WelcomeEmail:       "¡Bienvenido a nuestra aplicación!",
	SubjectKeysInstance.PasswordResetEmail: "Instrucciones para restablecer la contraseña",
	SubjectKeysInstance.OTPEmail:           "Su contraseña de un solo uso (OTP)",
}

type Subjects struct {
	En map[SubjectKeysEnum]string
	Es map[SubjectKeysEnum]string
}

var EmailSubjects = Subjects{
	En: EnSubjects,
	Es: EsSubjects,
}

func GetSubject(locale locales.LocaleTypeEnum, key SubjectKeysEnum) string {
	switch locale {
	case locales.EN_US:
		return EmailSubjects.En[key]
	case locales.ES_ES:
		return EmailSubjects.Es[key]
	default:
		return EmailSubjects.En[key]
	}
}
