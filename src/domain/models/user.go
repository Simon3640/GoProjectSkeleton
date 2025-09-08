package models

import "strconv"

type UserBase struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Phone    string `json:"phone"`
	Status   string `json:"status"`
	RoleID   uint   `json:"role_id"`
	OTPLogin bool   `json:"otp_login"`
}

func (u UserBase) Validate() []string {
	var errs []string

	if u.Name == "" {
		errs = append(errs, "name is required")
	}
	if u.Email == "" {
		errs = append(errs, "email is required")
	}
	if u.Phone == "" {
		errs = append(errs, "phone is required")
	}
	if u.Status == "" {
		errs = append(errs, "status is required")
	}
	if u.RoleID == 0 {
		errs = append(errs, "role_id is required")
	}
	// regex for email validation
	if !IsValidEmail(u.Email) {
		errs = append(errs, "email is invalid")
	}

	return errs
}

type UserWithRole struct {
	UserBase
	role Role
	ID   uint `json:"id"`
}

func (u *UserWithRole) SetRole(role Role) {
	u.role = role
}

func (u *UserWithRole) UserIsAdmin() bool {
	return u.role.Key == "admin"
}

func (u *UserWithRole) GetRoleKey() string {
	return u.role.Key
}

func (u *UserWithRole) GetUserID() uint {
	return u.ID
}

func (u *UserWithRole) GetUserIDString() string {
	return strconv.FormatUint(uint64(u.ID), 10)
}

type User struct {
	UserBase
	DBBaseModel
}

func (u User) GetUserID() uint {
	return u.ID
}

type UserInDB struct {
	UserBase
	DBBaseModel
}
