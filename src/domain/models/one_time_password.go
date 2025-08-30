package models

import "time"

type OneTimePasswordPurpose string

const (
	OneTimePasswordLogin OneTimePasswordPurpose = "login"
)

type OneTimePasswordBase struct {
	UserID  uint                   `json:"user_id"`
	Purpose OneTimePasswordPurpose `json:"purpose"`
	Hash    []byte                 `json:"hash"`
	IsUsed  bool                   `json:"is_used"`
	Expires time.Time              `json:"expires"`
}

func (o *OneTimePasswordBase) Validate() []string {
	var errs []string

	if o.UserID == 0 {
		errs = append(errs, "user_id is required")
	}
	if o.Purpose == "" {
		errs = append(errs, "purpose is required")
	}
	if len(o.Hash) == 0 {
		errs = append(errs, "hash is required")
	}
	if o.Expires.IsZero() {
		errs = append(errs, "expires is required")
	}

	if o.Purpose != OneTimePasswordLogin {
		errs = append(errs, "purpose is invalid")
	}

	return errs
}

type OneTimePassword struct {
	OneTimePasswordBase
	DBBaseModel
	Selector string `json:"selector"`
}
