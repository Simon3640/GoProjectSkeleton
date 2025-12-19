package passwordmocks

import (
	dtos "github.com/simon3640/goprojectskeleton/src/application/modules/password/dtos"
	"github.com/simon3640/goprojectskeleton/src/domain/models"
)

// PasswordBase is the base model for a password
var PasswordBase = models.PasswordBase{
	UserID:    1,
	Hash:      "$trongPassword123",
	ExpiresAt: nil,
	IsActive:  true,
}

// PasswordCreate is the create model for a password
var PasswordCreate = dtos.PasswordCreate{
	PasswordBase: PasswordBase,
}
