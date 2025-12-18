// Package passwordservices contains the services for the password module
package passwordservices

import (
	contractsProviders "github.com/simon3640/goprojectskeleton/src/application/contracts/providers"
	passwordcontracts "github.com/simon3640/goprojectskeleton/src/application/modules/password/contracts"
	dtos "github.com/simon3640/goprojectskeleton/src/application/modules/password/dtos"
	applicationerrors "github.com/simon3640/goprojectskeleton/src/application/shared/errors"
	"github.com/simon3640/goprojectskeleton/src/domain/models"
)

// CreatePasswordService creates a new password
func CreatePasswordService(
	passwordCreateNoHash dtos.PasswordCreateNoHash,
	hashProvider contractsProviders.IHashProvider,
	passwordRepository passwordcontracts.IPasswordRepository,
) (*models.Password, *applicationerrors.ApplicationError) {
	hashedPassword, err := hashProvider.HashPassword(passwordCreateNoHash.NoHashedPassword)
	if err != nil {
		return nil, err
	}
	passwordCreate := dtos.NewPasswordCreate(
		passwordCreateNoHash.UserID,
		hashedPassword,
		passwordCreateNoHash.ExpiresAt,
		passwordCreateNoHash.IsActive,
	)
	res, err := passwordRepository.Create(passwordCreate)
	if err != nil {
		return nil, err
	}
	return res, nil
}
