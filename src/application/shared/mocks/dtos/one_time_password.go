package dto_mocks

import (
	"encoding/hex"
	"time"

	"gormgoskeleton/src/domain/models"
)

var OneTimePassword = models.OneTimePassword{
	OneTimePasswordBase: models.OneTimePasswordBase{
		UserID:  1,
		Purpose: models.OneTimePasswordLogin,
		Hash:    []byte(hex.EncodeToString([]byte("hashed_otp"))),
		IsUsed:  false,
		Expires: time.Now().Add(10 * time.Minute),
	},
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
