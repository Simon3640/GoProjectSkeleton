// Package authservices contains the services for the auth module.
package authservices

import (
	contractproviders "github.com/simon3640/goprojectskeleton/src/application/contracts/providers"
	authcontracts "github.com/simon3640/goprojectskeleton/src/application/modules/auth/contracts"
	app_context "github.com/simon3640/goprojectskeleton/src/application/shared/context"
	"github.com/simon3640/goprojectskeleton/src/application/shared/locales"
	services "github.com/simon3640/goprojectskeleton/src/application/shared/services"
	emailservices "github.com/simon3640/goprojectskeleton/src/application/shared/services/emails"
	emailmodels "github.com/simon3640/goprojectskeleton/src/application/shared/services/emails/models"
	"github.com/simon3640/goprojectskeleton/src/application/shared/settings"
	"github.com/simon3640/goprojectskeleton/src/application/shared/templates"
	"github.com/simon3640/goprojectskeleton/src/domain/models"
)

// Verify that SendOTPEmailBackgroundService implements BackgroundService interface
var _ services.BackgroundService[SendOTPEmailInput] = (*SendOTPEmailBackgroundService)(nil)

// SendOTPEmailInput is the input for the SendOTPEmailBackgroundService
type SendOTPEmailInput struct {
	UserID   uint
	Email    string
	UserName string
}

// SendOTPEmailBackgroundService is a background service that creates an OTP and sends it via email
type SendOTPEmailBackgroundService struct {
	log          contractproviders.ILoggerProvider
	otpRepo      authcontracts.IOneTimePasswordRepository
	hashProvider contractproviders.IHashProvider
}

// NewSendOTPEmailBackgroundService creates a new instance of SendOTPEmailBackgroundService
func NewSendOTPEmailBackgroundService(
	log contractproviders.ILoggerProvider,
	otpRepo authcontracts.IOneTimePasswordRepository,
	hashProvider contractproviders.IHashProvider,
) *SendOTPEmailBackgroundService {
	return &SendOTPEmailBackgroundService{
		log:          log,
		otpRepo:      otpRepo,
		hashProvider: hashProvider,
	}
}

// Execute implements the BackgroundService interface
// It creates an OTP and sends it via email to the user
func (s *SendOTPEmailBackgroundService) Execute(
	_ *app_context.AppContext,
	locale locales.LocaleTypeEnum,
	input SendOTPEmailInput,
) error {
	// Create OTP
	otp, err := CreateOneTimePasswordService(
		input.UserID,
		models.OneTimePasswordLogin,
		s.hashProvider,
		s.otpRepo,
	)
	if err != nil {
		s.log.Error("Error creating OTP in background service", err.ToError())
		return err.ToError()
	}

	// Build email data
	otpEmailData := emailmodels.OneTimePasswordEmailData{
		Name:              input.UserName,
		OTPCode:           otp,
		ExpirationMinutes: int(settings.AppSettingsInstance.OneTimeTokenPasswordTTL),
		AppName:           settings.AppSettingsInstance.AppName,
		SupportEmail:      settings.AppSettingsInstance.AppSupportEmail,
	}

	// Send email
	if err := emailservices.OneTimePasswordEmailServiceInstance.SendWithTemplate(
		otpEmailData,
		input.Email,
		locale,
		templates.TemplateKeysInstance.OTPEmail,
		emailservices.SubjectKeysInstance.OTPEmail,
	); err != nil {
		s.log.Error("Error sending OTP email in background service", err.ToError())
		return err.ToError()
	}

	return nil
}

// Name returns the name of the service for logging and tracing
func (s *SendOTPEmailBackgroundService) Name() string {
	return "send-otp-email"
}
