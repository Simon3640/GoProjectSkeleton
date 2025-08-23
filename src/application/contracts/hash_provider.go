package contracts

import application_errors "gormgoskeleton/src/application/shared/errors"

type IHashProvider interface {
	HashPassword(password string) (string, *application_errors.ApplicationError)
	VerifyPassword(hashedPassword, password string) (bool, *application_errors.ApplicationError)
}
