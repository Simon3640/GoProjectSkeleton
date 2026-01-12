// Package authservices contains the services for the auth module.
package authservices

import (
	contractsProviders "github.com/simon3640/goprojectskeleton/src/application/contracts/providers"
	authcontracts "github.com/simon3640/goprojectskeleton/src/application/modules/auth/contracts"
	dtos "github.com/simon3640/goprojectskeleton/src/application/modules/auth/dtos"
	applicationerrors "github.com/simon3640/goprojectskeleton/src/application/shared/errors"
	sharedmodels "github.com/simon3640/goprojectskeleton/src/domain/shared/models"
)

func CreateOneTimePasswordService(
	userID uint,
	purpose sharedmodels.OneTimePasswordPurpose,
	hashProvider contractsProviders.IHashProvider,
	passwordRepository authcontracts.IOneTimePasswordRepository,
) (string, *applicationerrors.ApplicationError) {
	password, hash, err := hashProvider.GenerateOTP()
	if err != nil {
		return "", err
	}

	passwordCreate := dtos.NewOneTimePasswordCreate(userID, purpose, hash)
	_, err = passwordRepository.Create(*passwordCreate)
	if err != nil {
		return "", err
	}
	return password, nil
}
