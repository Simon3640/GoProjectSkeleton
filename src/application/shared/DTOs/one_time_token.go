package dtos

import (
	"time"

	"github.com/simon3640/goprojectskeleton/src/application/shared/settings"
	sharedmodels "github.com/simon3640/goprojectskeleton/src/domain/shared/models"
	usermodels "github.com/simon3640/goprojectskeleton/src/domain/user/models"
)

type OneTimeTokenCreate struct {
	sharedmodels.OneTimeTokenBase
}

// PurposeTokenToDuration converts a token purpose to its duration
func PurposeTokenToDuration(purpose sharedmodels.OneTimeTokenPurpose) time.Duration {
	switch purpose {
	case sharedmodels.OneTimeTokenPurposePasswordReset:
		return time.Duration(settings.AppSettingsInstance.OneTimeTokenPasswordTTL) * time.Minute
	case sharedmodels.OneTimeTokenPurposeEmailVerify:
		return time.Duration(settings.AppSettingsInstance.OneTimeTokenEmailVerifyTTL) * time.Minute
	default:
		return time.Duration(settings.AppSettingsInstance.OneTimeTokenEmailVerifyTTL) * time.Minute
	}
}

// NewOneTimeTokenCreate creates a new one-time token create DTO
func NewOneTimeTokenCreate(userID uint, purpose sharedmodels.OneTimeTokenPurpose, hash []byte) *OneTimeTokenCreate {
	// TODO: move expiration to another place
	return &OneTimeTokenCreate{
		OneTimeTokenBase: sharedmodels.OneTimeTokenBase{
			UserID:  userID,
			Purpose: purpose,
			Hash:    hash,
			Expires: time.Now().Add(PurposeTokenToDuration(purpose)),
			IsUsed:  false,
		},
	}
}

type OneTimeTokenUpdate struct {
	IsUsed bool `json:"isUsed,omitempty"`
	ID     uint `json:"id"`
}

type OneTimeTokenUser struct {
	User  usermodels.User
	Token string `json:"token"`
}

func (o *OneTimeTokenUser) BuildURL(baseURL string) string {
	return baseURL + "?token=" + o.Token
}
