package app_context

import (
	"context"

	"goprojectskeleton/src/domain/models"
)

type AppContext struct {
	context.Context
	User *models.User
}

func NewContextWithUser(user *models.User) *AppContext {
	return &AppContext{
		Context: context.Background(),
		User:    user,
	}
}
