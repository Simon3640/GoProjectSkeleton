package email_models

type NewUserEmailData struct {
	Name              string
	ActivationLink    string
	ExpirationMinutes int
	AppName           string
	SupportEmail      string
}
