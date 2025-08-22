package usecase

import "gormgoskeleton/src/domain/models"

type Guard func(user models.UserWithRole, input any) error

type Guards struct {
	list  []Guard
	actor models.UserWithRole
}

func (g Guards) Validate(input any) error {
	for _, guard := range g.list {
		if err := guard(g.actor, input); err != nil {
			return err
		}
	}
	return nil
}

func NewGuards(guards ...Guard) Guards {
	return Guards{
		list: guards,
	}
}

func (g *Guards) SetActor(actor models.UserWithRole) {
	g.actor = actor
}
