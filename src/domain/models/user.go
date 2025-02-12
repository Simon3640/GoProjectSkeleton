package models

type UserBase struct {
	Name   string `json:"name"`
	Email  string `json:"email"`
	Phone  string `json:"phone"`
	Status string `json:"status"`
}

type UserCreate struct {
	UserBase
}

type UserUpdate struct {
	Name   *string `json:"name"`
	Email  *string `json:"email"`
	Phone  *string `json:"phone"`
	Status *string `json:"status"`
}

type User struct {
	UserBase
	ID int `json:"id"`
}
