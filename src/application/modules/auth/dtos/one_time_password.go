package authdtos

import (
	"time"

	"github.com/simon3640/goprojectskeleton/src/application/shared/settings"
	sharedmodels "github.com/simon3640/goprojectskeleton/src/domain/shared/models"
)

// OneTimePasswordCreate is the DTO for creating a one time password
type OneTimePasswordCreate struct {
	sharedmodels.OneTimePasswordBase
}

// PurposePasswordToDuration converts the purpose to the duration
func PurposePasswordToDuration(purpose sharedmodels.OneTimePasswordPurpose) time.Duration {
	switch purpose {
	case sharedmodels.OneTimePasswordLogin:
		return time.Duration(settings.AppSettingsInstance.OneTimePasswordTTL) * time.Minute
	default:
		return time.Duration(settings.AppSettingsInstance.OneTimePasswordTTL) * time.Minute
	}
}

// NewOneTimePasswordCreate creates a new one time password create DTO
func NewOneTimePasswordCreate(userID uint, purpose sharedmodels.OneTimePasswordPurpose, hash []byte) *OneTimePasswordCreate {
	// TODO: move expiration to another place
	return &OneTimePasswordCreate{
		OneTimePasswordBase: sharedmodels.OneTimePasswordBase{
			UserID:  userID,
			Purpose: purpose,
			Hash:    hash,
			Expires: time.Now().Add(PurposePasswordToDuration(purpose)),
			IsUsed:  false,
		},
	}
}

// OneTimePasswordUpdate is the DTO for updating a one time password
type OneTimePasswordUpdate struct {
	IsUsed bool `json:"isUsed,omitempty"`
	ID     uint `json:"id"`
}
