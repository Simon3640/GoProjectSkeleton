package dtos

import (
	"time"

	"github.com/simon3640/goprojectskeleton/src/application/shared/settings"
	"github.com/simon3640/goprojectskeleton/src/domain/models"
)

type OneTimeTokenCreate struct {
	models.OneTimeTokenBase
}

func PurposeTokenToDuration(purpose models.OneTimeTokenPurpose) time.Duration {
	switch purpose {
	case models.OneTimeTokenPurposePasswordReset:
		return time.Duration(settings.AppSettingsInstance.OneTimeTokenPasswordTTL) * time.Minute
	case models.OneTimeTokenPurposeEmailVerify:
		return time.Duration(settings.AppSettingsInstance.OneTimeTokenEmailVerifyTTL) * time.Minute
	default:
		return time.Duration(settings.AppSettingsInstance.OneTimeTokenEmailVerifyTTL) * time.Minute
	}
}

func NewOneTimeTokenCreate(userID uint, purpose models.OneTimeTokenPurpose, hash []byte) *OneTimeTokenCreate {
	// TODO: move expiration to another place
	return &OneTimeTokenCreate{
		OneTimeTokenBase: models.OneTimeTokenBase{
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
	User  models.User
	Token string `json:"token"`
}

func (o *OneTimeTokenUser) BuildURL(baseURL string) string {
	return baseURL + "?token=" + o.Token
}
