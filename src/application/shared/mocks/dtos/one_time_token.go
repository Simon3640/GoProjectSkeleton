package dtomocks

import (
	"encoding/hex"
	dtos "goprojectskeleton/src/application/shared/DTOs"
	"goprojectskeleton/src/domain/models"
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
