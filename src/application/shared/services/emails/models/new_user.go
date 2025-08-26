package email_models

import "time"

type NewUserEmailData struct {
	Name            string
	ActivationToken string
	Expiration      time.Time
	AppName         string
	SupportEmail    string
}
