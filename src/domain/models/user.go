package models

type UserBase struct {
	Name   string `json:"name"`
	Email  string `json:"email"`
	Phone  string `json:"phone"`
	Status string `json:"status"`
	RoleID uint   `json:"role_id"`
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

type UserCreate struct {
	UserBase
}

// cant create admin role
func (u UserCreate) Validate() []string {
	// super Validate
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

func (u UserAndPasswordCreate) Validate() []string {
	errs := u.UserCreate.Validate()
	if !IsValidPassword(u.Password) {
		errs = append(errs, "password is invalid")
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
	if u.Email != nil && *u.Email == "" && !IsValidEmail(*u.Email) {
		errs = append(errs, "email is invalid")
	}
	return errs
}

func (u UserUpdate) GetUserID() uint {
	return u.ID
}

type User struct {
	UserBase
	ID uint `json:"id"`
}

func (u User) GetUserID() uint {
	return u.ID
}

type UserInDB struct {
	UserBase
	DBBaseModel
}
