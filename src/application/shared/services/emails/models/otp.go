package email_models

type OneTimePasswordEmailData struct {
	Name              string
	OTPCode           string
	ExpirationMinutes int
	AppName           string
	SupportEmail      string
}
