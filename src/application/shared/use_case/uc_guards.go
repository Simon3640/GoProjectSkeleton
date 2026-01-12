package usecase

import (
	"github.com/simon3640/goprojectskeleton/src/application/shared/locales/messages"
	usermodels "github.com/simon3640/goprojectskeleton/src/domain/user/models"
)

// Guard is a function that validates user access to a resource
type Guard func(user usermodels.UserWithRole, input any) *messages.MessageKeysEnum

type Guards struct {
	list  []Guard
	actor usermodels.UserWithRole
}

// Validate validates the input against the guards
func (g Guards) Validate(input any) *messages.MessageKeysEnum {
	for _, guard := range g.list {
		if err := guard(g.actor, input); err != nil {
			return err
		}
	}
	return nil
}

// NewGuards creates a new Guards instance
func NewGuards(guards ...Guard) Guards {
	return Guards{
		list: guards,
	}
}

// SetActor sets the actor for the guards
func (g *Guards) SetActor(actor usermodels.UserWithRole) {
	g.actor = actor
}
