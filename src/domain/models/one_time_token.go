package models

import "time"

type OneTimeTokenBase struct {
	UserID  uint      `json:"user_id"`
	Purpose string    `json:"purpose"`
	Hash    string    `json:"hash"`
	IsUsed  bool      `json:"is_used"`
	Expires time.Time `json:"expires"`
}

type OneTimeToken struct {
	OneTimeTokenBase
	DBBaseModel
}
