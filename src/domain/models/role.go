package models

type RoleBase struct {
	Key      string `json:"key"`
	IsActive bool   `json:"status"`
	Priority int    `json:"priority"`
}

type RoleCreate struct {
	RoleBase
}

type RoleUpdateBase struct {
	Key      *string `json:"key"`
	IsActive *bool   `json:"status"`
	Priority *int    `json:"priority"`
}

type RoleUpdate struct {
	RoleUpdateBase
	ID uint `json:"id"`
}

type Role struct {
	RoleBase
	ID uint `json:"id"`
}

type RoleInDB struct {
	RoleBase
	DBBaseModel
}
