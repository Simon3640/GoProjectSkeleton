package models

import "strconv"

// UserStatus is the status of a user
// It can be:
// - pending
// - active
// - inactive
// - suspended
// - deleted
type UserStatus string

const (
	// UserStatusPending is the status of a user when they are pending
	UserStatusPending UserStatus = "pending"
	// UserStatusActive is the status of a user when they are active
	UserStatusActive UserStatus = "active"
	// UserStatusInactive is the status of a user when they are inactive
	UserStatusInactive UserStatus = "inactive"
	// UserStatusSuspended is the status of a user when they are suspended
	UserStatusSuspended UserStatus = "suspended"
	// UserStatusDeleted is the status of a user when they are deleted
	UserStatusDeleted UserStatus = "deleted"
)

func (s UserStatus) String() string {
	return map[UserStatus]string{
		UserStatusPending:   "pending",
		UserStatusActive:    "active",
		UserStatusInactive:  "inactive",
		UserStatusSuspended: "suspended",
		UserStatusDeleted:   "deleted",
	}[s]
}

type UserBase struct {
	Name     string      `json:"name"`
	Email    string      `json:"email"`
	Phone    string      `json:"phone"`
	Status   *UserStatus `json:"status,omitempty"`
	RoleID   uint        `json:"role_id"`
	OTPLogin bool        `json:"otpLogin"`
}

// Validate validates the user base
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
	if u.RoleID == 0 {
		errs = append(errs, "role_id is required")
	}
	// regex for email validation
	if !IsValidEmail(u.Email) {
		errs = append(errs, "email is invalid")
	}

	return errs
}

// ValidateCreate validates the user base for creation
func (u *UserBase) ValidateCreate() []string {
	errs := u.Validate()
	status := UserStatusPending
	u.Status = &status
	if u.RoleID == 1 { // TODO: replace with constant
		errs = append(errs, "admin role is not allowed")
	}
	return errs
}

// UserUpdateBase is the update base for a user
type UserUpdateBase struct {
	Name     *string     `json:"name"`
	Email    *string     `json:"email"`
	Phone    *string     `json:"phone"`
	Status   *UserStatus `json:"status,omitempty"`
	RoleID   *uint       `json:"role_id,omitempty"`
	OTPLogin *bool       `json:"otpLogin,omitempty"`
}

// Validate validates the user update base
func (u *UserUpdateBase) Validate() []string {
	var errs []string
	if u.Email != nil && !IsValidEmail(*u.Email) {
		errs = append(errs, "email is invalid")
	}
	if u.RoleID != nil && *u.RoleID == 1 { // TODO: replace with constant
		errs = append(errs, "admin role is not allowed")
	}
	return errs
}

// UserUpdate is the update structure for a user
type UserUpdate struct {
	UserUpdateBase
	ID uint `json:"id"`
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
