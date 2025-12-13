package app_context

import (
	"context"

	dtos "github.com/simon3640/goprojectskeleton/src/application/shared/DTOs"
	"github.com/simon3640/goprojectskeleton/src/domain/models"
)

// AppContext is the context for the application layer
type AppContext struct {
	context.Context
	User         *models.UserWithRole
	OneTimeToken *dtos.OneTimeTokenUser
}

// NewContextWithUser creates a new AppContext with a user
func NewContextWithUser(user *models.UserWithRole) *AppContext {
	return &AppContext{
		Context: context.Background(),
		User:    user,
	}
}

// NewVoidAppContext creates a new AppContext with a void context
func NewVoidAppContext() *AppContext {
	return &AppContext{
		Context: context.Background(),
	}
}

// AddUserToContext adds a user to the context
func (a *AppContext) AddUserToContext(user *models.UserWithRole) {
	a.User = user
}

// AddOneTimeTokenToContext adds a one-time token to the context
func (a *AppContext) AddOneTimeTokenToContext(oneTimeToken dtos.OneTimeTokenUser) {
	a.OneTimeToken = &oneTimeToken
}
