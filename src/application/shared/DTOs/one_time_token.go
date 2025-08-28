package dtos

import (
	"gormgoskeleton/src/application/shared/settings"
	"gormgoskeleton/src/domain/models"
	"time"
)

type OneTimeTokenCreate struct {
	models.OneTimeTokenBase
}

func NewOneTimeTokenCreate(userID uint, purpose string, hash []byte) *OneTimeTokenCreate {
	// TODO: move expiration to another place
	expires := time.Now().Add(time.Duration(settings.AppSettingsInstance.OneTimeTokenTTL) * time.Minute)
	return &OneTimeTokenCreate{
		OneTimeTokenBase: models.OneTimeTokenBase{
			UserID:  userID,
			Purpose: purpose,
			Hash:    hash,
			Expires: expires,
			IsUsed:  false,
		},
	}
}

type OneTimeTokenUpdate struct {
	IsUsed *bool `json:"is_used,omitempty"`
	ID     uint  `json:"id"`
}
