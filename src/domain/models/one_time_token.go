package models

import "time"

type OneTimeTokenPurpose string

const (
	OneTimeTokenPurposePasswordReset OneTimeTokenPurpose = "password_reset"
	OneTimeTokenPurposeEmailVerify   OneTimeTokenPurpose = "email_verify"
)

type OneTimeTokenBase struct {
	UserID  uint                `json:"user_id"`
	Purpose OneTimeTokenPurpose `json:"purpose"`
	Hash    []byte              `json:"hash"`
	IsUsed  bool                `json:"is_used"`
	Expires time.Time           `json:"expires"`
}

func (o *OneTimeTokenBase) Validate() []string {
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

	if o.Purpose != OneTimeTokenPurposePasswordReset && o.Purpose != OneTimeTokenPurposeEmailVerify {
		errs = append(errs, "purpose is invalid")
	}

	return errs
}

type OneTimeToken struct {
	OneTimeTokenBase
	DBBaseModel
}
