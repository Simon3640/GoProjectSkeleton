package services

import (
	contractsProviders "github.com/simon3640/goprojectskeleton/src/application/contracts/providers"
	contracts_repositories "github.com/simon3640/goprojectskeleton/src/application/contracts/repositories"
	dtos "github.com/simon3640/goprojectskeleton/src/application/shared/DTOs"
	application_errors "github.com/simon3640/goprojectskeleton/src/application/shared/errors"
	"github.com/simon3640/goprojectskeleton/src/domain/models"
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
