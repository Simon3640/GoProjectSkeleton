package dtomocks

import (
	"encoding/hex"

	dtos "github.com/simon3640/goprojectskeleton/src/application/shared/DTOs"
	sharedmodels "github.com/simon3640/goprojectskeleton/src/domain/shared/models"
)

// OneTimeTokenBase is a mock one-time token base for testing
var OneTimeTokenBase = sharedmodels.OneTimeTokenBase{
	UserID:  1,
	Purpose: sharedmodels.OneTimeTokenPurposeEmailVerify,
	Hash:    []byte(hex.EncodeToString([]byte("hashed_ott"))),
	IsUsed:  false,
}

var OneTimeTokenCreate = dtos.OneTimeTokenCreate{
	OneTimeTokenBase: OneTimeTokenBase,
}
