package guards

import (
	"github.com/simon3640/goprojectskeleton/src/application/shared/locales/messages"
	usecase "github.com/simon3640/goprojectskeleton/src/application/shared/use_case"
	sharedmodels "github.com/simon3640/goprojectskeleton/src/domain/shared/models"
	usermodels "github.com/simon3640/goprojectskeleton/src/domain/user/models"
)

func RoleGuard(allowedRoles ...string) usecase.Guard {
	return func(user usermodels.UserWithRole, _ any) *messages.MessageKeysEnum {
		for _, role := range allowedRoles {
			if user.GetRoleKey() == role {
				return nil
			}
		}
		return &messages.MessageKeysInstance.UNAUTHORIZED_RESOURCE
	}
}

// Partial Resource with UserID
func UserResourceGuard[T sharedmodels.HasUserID]() usecase.Guard {
	return func(user usermodels.UserWithRole, input any) *messages.MessageKeysEnum {
		resource, ok := input.(T)
		if !ok {
			return &messages.MessageKeysInstance.UNAUTHORIZED_RESOURCE
		}
		if user.ID != resource.GetUserID() {
			return &messages.MessageKeysInstance.UNAUTHORIZED_RESOURCE
		}
		return nil
	}
}

// UserGetItSelf checks if the user is trying to get their own data
func UserGetItSelf(user usermodels.UserWithRole, input any) *messages.MessageKeysEnum {
	id, ok := input.(uint)
	if !ok {
		return &messages.MessageKeysInstance.SOMETHING_WENT_WRONG
	}
	if user.ID != id {
		return &messages.MessageKeysInstance.UNAUTHORIZED_RESOURCE
	}
	return nil
}
