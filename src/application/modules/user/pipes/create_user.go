// Package userpipes provides pipes for user use cases
package userpipes

import (
	"context"

	userdtos "github.com/simon3640/goprojectskeleton/src/application/modules/user/dtos"
	userusecases "github.com/simon3640/goprojectskeleton/src/application/modules/user/use_cases"
	"github.com/simon3640/goprojectskeleton/src/application/shared/locales"
	usecase "github.com/simon3640/goprojectskeleton/src/application/shared/use_case"
	workers "github.com/simon3640/goprojectskeleton/src/application/shared/workers"
	"github.com/simon3640/goprojectskeleton/src/domain/models"
)

// NewCreateUserPipe creates a new create user pipe.
// The user is created synchronously, and the welcome email is sent in background.
// This allows returning the response immediately without waiting for the email to be sent.
func NewCreateUserPipe(
	ctx context.Context,
	locale locales.LocaleTypeEnum,
	createUserPasswordUC *userusecases.CreateUserAndPasswordUseCase,
	createUserSendEmailUseCase *userusecases.CreateUserSendEmailUseCase,
) *usecase.DAG[userdtos.UserAndPasswordCreate, models.User] {
	dag := usecase.NewDag(ctx, usecase.NewStep(createUserPasswordUC), locale, workers.GetBackgroundExecutor())
	// The email is sent in background, the response is returned immediately
	return usecase.ThenBackground(dag, usecase.NewStep(createUserSendEmailUseCase), "create-user-send-email")
}
