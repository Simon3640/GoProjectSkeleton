package app_context

import (
	"context"

	dtos "github.com/simon3640/goprojectskeleton/src/application/shared/DTOs"
	"github.com/simon3640/goprojectskeleton/src/domain/models"
)

type AppContext struct {
	context.Context
	User         *models.UserWithRole
	OneTimeToken *dtos.OneTimeTokenUser
}

func NewContextWithUser(user *models.UserWithRole) *AppContext {
	return &AppContext{
		Context: context.Background(),
		User:    user,
	}
}

func NewVoidAppContext() *AppContext {
	return &AppContext{
		Context: context.Background(),
	}
}

func (a *AppContext) AddUserToContext(user *models.UserWithRole) {
	a.User = user
}

func (a *AppContext) AddOneTimeTokenToContext(oneTimeToken dtos.OneTimeTokenUser) {
	a.OneTimeToken = &oneTimeToken
}
