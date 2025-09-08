package services

import (
	contracts_providers "gormgoskeleton/src/application/contracts/providers"
	contracts_repositories "gormgoskeleton/src/application/contracts/repositories"
	dtos "gormgoskeleton/src/application/shared/DTOs"
	application_errors "gormgoskeleton/src/application/shared/errors"
	"gormgoskeleton/src/domain/models"
)

func CreateOneTimeTokenService(
	userID uint,
	purpose models.OneTimeTokenPurpose,
	hashProvider contracts_providers.IHashProvider,
	tokenRepository contracts_repositories.IOneTimeTokenRepository,
) (string, *application_errors.ApplicationError) {
	token, hash, err := hashProvider.OneTimeToken()
	if err != nil {
		return "", err
	}

	tokenCreate := dtos.NewOneTimeTokenCreate(userID, purpose, hash)
	_, err = tokenRepository.Create(*tokenCreate)
	if err != nil {
		return "", err
	}
	return token, nil
}
