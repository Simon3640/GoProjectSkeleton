package dtos

import "gormgoskeleton/src/domain/models"

type OneTimeTokenCreate struct {
	models.OneTimeTokenBase
}

type OneTimeTokenUpdate struct {
	IsUsed *bool `json:"is_used,omitempty"`
	ID     uint  `json:"id"`
}
