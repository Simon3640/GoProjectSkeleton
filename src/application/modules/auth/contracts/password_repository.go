package authcontracts

import (
	applicationerrors "github.com/simon3640/goprojectskeleton/src/application/shared/errors"
	passwordmodels "github.com/simon3640/goprojectskeleton/src/domain/password/models"
)

// IPasswordRepository is the interface for the password repository
type IPasswordRepository interface {
	// GetActivePassword gets the active password for a user
	GetActivePassword(userEmail string) (*passwordmodels.Password, *applicationerrors.ApplicationError)
}
