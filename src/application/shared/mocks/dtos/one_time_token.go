package dtomocks

import (
	"encoding/hex"

	dtos "github.com/simon3640/goprojectskeleton/src/application/shared/DTOs"
	"github.com/simon3640/goprojectskeleton/src/domain/models"
)

var OneTimeTokenBase = models.OneTimeTokenBase{
	UserID:  1,
	Purpose: models.OneTimeTokenPurposeEmailVerify,
	Hash:    []byte(hex.EncodeToString([]byte("hashed_ott"))),
	IsUsed:  false,
}

var OneTimeTokenCreate = dtos.OneTimeTokenCreate{
	OneTimeTokenBase: OneTimeTokenBase,
}
