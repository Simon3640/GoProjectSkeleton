package dtos

import "gormgoskeleton/src/domain/models"

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

type UserMultiResponse = MultipleResponse[models.User]
type UserSingleResponse = SingleResponse[models.User]
