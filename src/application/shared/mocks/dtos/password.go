package dtomocks

import (
	dtos "gormgoskeleton/src/application/shared/DTOs"
	"gormgoskeleton/src/domain/models"
)

var PasswordBase = models.PasswordBase{
	UserID:    1,
	Hash:      "$trongPassword123",
	ExpiresAt: nil,
	IsActive:  true,
}

var PasswordCreate = dtos.PasswordCreate{
	PasswordBase: PasswordBase,
}
