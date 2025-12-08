package dtos

import (
	"time"

	"github.com/simon3640/goprojectskeleton/src/application/shared/settings"
	"github.com/simon3640/goprojectskeleton/src/domain/models"
)

type OneTimePasswordCreate struct {
	models.OneTimePasswordBase
}

func PurposePasswordToDuration(purpose models.OneTimePasswordPurpose) time.Duration {
	switch purpose {
	case models.OneTimePasswordLogin:
		return time.Duration(settings.AppSettingsInstance.OneTimePasswordTTL) * time.Minute
	default:
		return time.Duration(settings.AppSettingsInstance.OneTimePasswordTTL) * time.Minute
	}
}

func NewOneTimePasswordCreate(userID uint, purpose models.OneTimePasswordPurpose, hash []byte) *OneTimePasswordCreate {
	// TODO: move expiration to another place
	return &OneTimePasswordCreate{
		OneTimePasswordBase: models.OneTimePasswordBase{
			UserID:  userID,
			Purpose: purpose,
			Hash:    hash,
			Expires: time.Now().Add(PurposePasswordToDuration(purpose)),
			IsUsed:  false,
		},
	}
}

type OneTimePasswordUpdate struct {
	IsUsed bool `json:"isUsed,omitempty"`
	ID     uint `json:"id"`
}
