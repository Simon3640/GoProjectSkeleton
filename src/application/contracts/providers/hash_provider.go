package contracts_providers

import application_errors "gormgoskeleton/src/application/shared/errors"

type IHashProvider interface {
	HashPassword(password string) (string, *application_errors.ApplicationError)
	VerifyPassword(hashedPassword, password string) (bool, *application_errors.ApplicationError)
	OneTimeToken() (string, []byte, *application_errors.ApplicationError)
	HashOneTimeToken(token string) []byte
	ValidateOneTimeToken(hashedToken []byte, token string) bool
}
