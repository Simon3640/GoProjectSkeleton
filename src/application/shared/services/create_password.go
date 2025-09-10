package services

import (
	contractsProviders "gormgoskeleton/src/application/contracts/providers"
	contracts_repositories "gormgoskeleton/src/application/contracts/repositories"
	dtos "gormgoskeleton/src/application/shared/DTOs"
	application_errors "gormgoskeleton/src/application/shared/errors"
	"gormgoskeleton/src/domain/models"
)

func CreatePasswordService(
	passwordCreateNoHash dtos.PasswordCreateNoHash,
	hashProvider contractsProviders.IHashProvider,
	passwordRepository contracts_repositories.IPasswordRepository,
) (*models.Password, *application_errors.ApplicationError) {
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
