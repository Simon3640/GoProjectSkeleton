package dtos

import "github.com/simon3640/goprojectskeleton/src/domain/models"

type UserCreate struct {
	models.UserBase
}

// cant create admin role
func (u *UserCreate) Validate() []string {
	return u.ValidateCreate()
}

type UserAndPasswordCreate struct {
	UserCreate
	Password string `json:"password"`
}

func (u *UserAndPasswordCreate) Validate() []string {
	errs := u.UserCreate.Validate()
	if !models.IsValidPassword(u.Password) {
		errs = append(errs, "password is invalid")
	}
	return errs
}

type UserUpdate struct {
	models.UserUpdateBase
	ID uint `json:"id"`
}

func (u UserUpdate) Validate() []string {
	return u.UserUpdateBase.Validate()
}

func (u UserUpdate) GetUserID() uint {
	return u.ID
}

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

type UserMultiResponse = MultipleResponse[models.User]
type UserSingleResponse = SingleResponse[models.User]
