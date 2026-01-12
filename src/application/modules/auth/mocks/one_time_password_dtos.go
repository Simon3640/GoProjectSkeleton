package authmocks

import (
	"encoding/hex"
	"time"

	dtos "github.com/simon3640/goprojectskeleton/src/application/modules/auth/dtos"
	sharedmodels "github.com/simon3640/goprojectskeleton/src/domain/shared/models"
)

// OneTimePasswordBase is a mock one-time password base for testing
var OneTimePasswordBase = sharedmodels.OneTimePasswordBase{
	UserID:  1,
	Purpose: sharedmodels.OneTimePasswordLogin,
	Hash:    []byte(hex.EncodeToString([]byte("hashed_otp"))),
	IsUsed:  false,
	Expires: time.Now().Add(10 * time.Minute),
}

var OneTimePasswordCreate = dtos.OneTimePasswordCreate{
	OneTimePasswordBase: OneTimePasswordBase,
}

// OneTimePassword is a mock one-time password for testing
var OneTimePassword = sharedmodels.OneTimePassword{
	OneTimePasswordBase: OneTimePasswordBase,
	DBBaseModel: sharedmodels.DBBaseModel{
		ID:        1,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	},
}

// ExpiredOneTimePassword is a mock expired one-time password for testing
var ExpiredOneTimePassword = sharedmodels.OneTimePassword{
	OneTimePasswordBase: sharedmodels.OneTimePasswordBase{
		UserID:  1,
		Purpose: sharedmodels.OneTimePasswordLogin,
		Hash:    []byte(hex.EncodeToString([]byte("hashed_otp"))),
		IsUsed:  false,
		Expires: time.Now().Add(-10 * time.Minute),
	},
	DBBaseModel: sharedmodels.DBBaseModel{
		ID:        1,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	},
}
