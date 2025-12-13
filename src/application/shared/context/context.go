package app_context

import (
	"context"

	"github.com/simon3640/goprojectskeleton/src/domain/models"
)

type AppContext struct {
	context.Context
	User *models.UserWithRole
}

func NewContextWithUser(user *models.UserWithRole) *AppContext {
	return &AppContext{
		Context: context.Background(),
		User:    user,
	}
}

func (a *AppContext) AddUserToContext(user *models.UserWithRole) {
	a.User = user
}
