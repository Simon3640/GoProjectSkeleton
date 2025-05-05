package models

type RoleBase struct {
	Key      string `json:"key"`
	IsActive string `json:"status"`
	Priority int    `json:"priority"`
}

type RoleCreate struct {
	RoleBase
}

type RoleUpdateBase struct {
	Key      *string `json:"key"`
	IsActive *string `json:"status"`
	Priority *int    `json:"priority"`
}

type RoleUpdate struct {
	RoleUpdateBase
	ID int `json:"id"`
}

type Role struct {
	RoleBase
	ID int `json:"id"`
}
