package authcontracts

import (
	application_errors "github.com/simon3640/goprojectskeleton/src/application/shared/errors"
	"github.com/simon3640/goprojectskeleton/src/domain/models"
)

// IPasswordRepository is the interface for the password repository
type IPasswordRepository interface {
	// GetActivePassword gets the active password for a user
	GetActivePassword(userEmail string) (*models.Password, *application_errors.ApplicationError)
}
