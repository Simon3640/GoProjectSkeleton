// Package userdtos contains the DTOs for the user module
package userdtos

import (
	shareddtos "github.com/simon3640/goprojectskeleton/src/application/shared/DTOs"
	"github.com/simon3640/goprojectskeleton/src/domain/models"
)

// UserCreate is the create structure for a user
type UserCreate struct {
	models.UserBase
}

// Validate validates the user create
// This validates also create defaults that are set in the ValidateCreate method, this is a pointer method because it needs to modify the user base
func (u *UserCreate) Validate() []string {
	return u.ValidateCreate()
}

// UserAndPasswordCreate is the create structure for a user and password
type UserAndPasswordCreate struct {
	UserCreate
	Password string `json:"password"`
}

// Validate validates the user and password create
func (u *UserAndPasswordCreate) Validate() []string {
	errs := u.UserCreate.Validate()
	if !models.IsValidPassword(u.Password) {
		errs = append(errs, "password is invalid")
	}
	return errs
}

// UserUpdate is the update structure for a user
type UserUpdate struct {
	models.UserUpdateBase
	ID uint `json:"id"`
}

// Validate validates the user update
func (u UserUpdate) Validate() []string {
	return u.UserUpdateBase.Validate()
}

// GetUserID returns the user ID
func (u UserUpdate) GetUserID() uint {
	return u.ID
}

// UserActivate is the activate structure for a user
type UserActivate struct {
	Token string `json:"token"`
}

// ResendWelcomeEmailRequest is the request for the resend welcome email use case
type ResendWelcomeEmailRequest struct {
	Email string `json:"email"`
}

// Validate validates the resend welcome email request
// returns a list of errors if the request is invalid
func (r *ResendWelcomeEmailRequest) Validate() []string {
	var errs []string
	if r.Email == "" {
		errs = append(errs, "email is required")
	}
	if !models.IsValidEmail(r.Email) {
		errs = append(errs, "email is invalid")
	}
	return errs
}

// UserMultiResponse is the multiple response for a user
type UserMultiResponse = shareddtos.MultipleResponse[models.User]

// UserSingleResponse is the single response for a user
type UserSingleResponse = shareddtos.SingleResponse[models.User]
