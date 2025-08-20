package models

type UserBase struct {
	Name   string `json:"name"`
	Email  string `json:"email"`
	Phone  string `json:"phone"`
	Status string `json:"status"`
	RoleID uint   `json:"role_id"`
}

type UserCreate struct {
	UserBase
}

type UserAndPasswordCreate struct {
	UserCreate
	Password string `json:"password"`
}

type UserRole struct {
	UserBase
	Role Role `json:"role"`
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

type User struct {
	UserBase
	ID uint `json:"id"`
}

type UserInDB struct {
	UserBase
	DBBaseModel
}
