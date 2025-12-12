// Package userpipes provides pipes for user use cases
package userpipes

import (
	"context"

	userdtos "github.com/simon3640/goprojectskeleton/src/application/modules/user/dtos"
	userusecases "github.com/simon3640/goprojectskeleton/src/application/modules/user/use_cases"
	"github.com/simon3640/goprojectskeleton/src/application/shared/locales"
	usecase "github.com/simon3640/goprojectskeleton/src/application/shared/use_case"
	"github.com/simon3640/goprojectskeleton/src/domain/models"
)

// NewCreateUserPipe creates a new create user pipe.
// El usuario se crea de forma síncrona, y el email de bienvenida se envía en background.
// Esto permite retornar la respuesta inmediatamente sin esperar el envío del email.
func NewCreateUserPipe(
	ctx context.Context,
	locale locales.LocaleTypeEnum,
	createUserPasswordUC *userusecases.CreateUserAndPasswordUseCase,
	createUserSendEmailUseCase *userusecases.CreateUserSendEmailUseCase,
) *usecase.DAG[userdtos.UserAndPasswordCreate, models.User] {
	dag := usecase.NewDag(usecase.NewStep(createUserPasswordUC), locale, ctx)
	// El email se envía en background, la respuesta se retorna inmediatamente
	return usecase.ThenBackground(dag, usecase.NewStep(createUserSendEmailUseCase))
}
