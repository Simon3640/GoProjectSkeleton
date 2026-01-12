// Package passwordcontracts contains the interfaces for the password repository
package passwordcontracts

import (
	contractsrepositories "github.com/simon3640/goprojectskeleton/src/application/contracts/repositories"
	dtos "github.com/simon3640/goprojectskeleton/src/application/modules/password/dtos"
	applicationerrors "github.com/simon3640/goprojectskeleton/src/application/shared/errors"
	passwordmodels "github.com/simon3640/goprojectskeleton/src/domain/password/models"
)

// IPasswordRepository is the interface for the password repository
type IPasswordRepository interface {
	contractsrepositories.IRepositoryBase[dtos.PasswordCreate, dtos.PasswordUpdate, passwordmodels.Password, passwordmodels.PasswordInDB]
	GetActivePassword(userEmail string) (*passwordmodels.Password, *applicationerrors.ApplicationError)
}
