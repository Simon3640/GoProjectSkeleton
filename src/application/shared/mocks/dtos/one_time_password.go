package dtomocks

import (
	"encoding/hex"
	"time"

	dtos "goprojectskeleton/src/application/shared/DTOs"
	"goprojectskeleton/src/domain/models"
)

var OneTimePasswordBase = models.OneTimePasswordBase{
	UserID:  1,
	Purpose: models.OneTimePasswordLogin,
	Hash:    []byte(hex.EncodeToString([]byte("hashed_otp"))),
	IsUsed:  false,
	Expires: time.Now().Add(10 * time.Minute),
}

var OneTimePasswordCreate = dtos.OneTimePasswordCreate{
	OneTimePasswordBase: OneTimePasswordBase,
}

var OneTimePassword = models.OneTimePassword{
	OneTimePasswordBase: OneTimePasswordBase,
	DBBaseModel: models.DBBaseModel{
		ID:        1,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	},
}

var ExpiredOneTimePassword = models.OneTimePassword{
	OneTimePasswordBase: models.OneTimePasswordBase{
		UserID:  1,
		Purpose: models.OneTimePasswordLogin,
		Hash:    []byte(hex.EncodeToString([]byte("hashed_otp"))),
		IsUsed:  false,
		Expires: time.Now().Add(-10 * time.Minute),
	},
	DBBaseModel: models.DBBaseModel{
		ID:        1,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	},
}
