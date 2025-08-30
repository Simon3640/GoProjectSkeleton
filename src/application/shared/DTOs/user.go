package dtos

import "gormgoskeleton/src/domain/models"

type UserCreate struct {
	models.UserBase
}

// cant create admin role
func (u *UserCreate) Validate() []string {
	// super Validate
	u.Status = "pending" // default status
	errs := u.UserBase.Validate()
	if u.RoleID == 1 { // TODO: replace with constant
		errs = append(errs, "admin role is not allowed")
	}
	return errs
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

type UserUpdateBase struct {
	Name   *string `json:"name"`
	Email  *string `json:"email"`
	Phone  *string `json:"phone"`
	Status *string `json:"status"`
	RoleID *uint   `json:"role_id,omitempty"`
}

type UserUpdate struct {
	UserUpdateBase
	ID uint `json:"id"`
}

func (u UserUpdate) Validate() []string {
	var errs []string
	if u.Email != nil && !models.IsValidEmail(*u.Email) {
		errs = append(errs, "email is invalid")
	}
	if u.RoleID != nil && *u.RoleID == 1 { // TODO: replace with constant
		errs = append(errs, "admin role is not allowed")
	}
	return errs
}

func (u UserUpdate) GetUserID() uint {
	return u.ID
}

type UserMultiResponse = MultipleResponse[models.User]
type UserSingleResponse = SingleResponse[models.User]
