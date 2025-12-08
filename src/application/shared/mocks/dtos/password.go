package dtomocks

import (
	dtos "github.com/simon3640/goprojectskeleton/src/application/shared/DTOs"
	"github.com/simon3640/goprojectskeleton/src/domain/models"
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
