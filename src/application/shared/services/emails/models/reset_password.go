package email_models

type ResetPasswordEmailData struct {
	Name              string
	ResetLink         string
	ExpirationMinutes int64
	AppName           string
	SupportEmail      string
}
