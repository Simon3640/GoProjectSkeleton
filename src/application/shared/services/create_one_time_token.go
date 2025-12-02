package services

import (
	contractsProviders "goprojectskeleton/src/application/contracts/providers"
	contracts_repositories "goprojectskeleton/src/application/contracts/repositories"
	dtos "goprojectskeleton/src/application/shared/DTOs"
	application_errors "goprojectskeleton/src/application/shared/errors"
	"goprojectskeleton/src/domain/models"
)

func CreateOneTimeTokenService(
	userID uint,
	purpose models.OneTimeTokenPurpose,
	hashProvider contractsProviders.IHashProvider,
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
