package guards

import (
	"gormgoskeleton/src/application/shared/locales/messages"
	usecase "gormgoskeleton/src/application/shared/use_case"
	"gormgoskeleton/src/domain/models"
)

func RoleGuard(allowedRoles ...string) usecase.Guard {
	return func(user models.UserWithRole, input any) *messages.MessageKeysEnum {
		for _, role := range allowedRoles {
			if user.GetRoleKey() == role {
				return nil
			}
		}
		return &messages.MessageKeysInstance.UNAUTHORIZED_RESOURCE
	}
}

// Partial Resource with UserID
func UserResourceGuard[T models.HasUserID]() usecase.Guard {
	return func(user models.UserWithRole, input any) *messages.MessageKeysEnum {
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

func UserGetItSelf(user models.UserWithRole, input any) *messages.MessageKeysEnum {
	id, ok := input.(uint)
	if !ok {
		return &messages.MessageKeysInstance.SOMETHING_WENT_WRONG
	}
	if user.ID != id {
		return &messages.MessageKeysInstance.UNAUTHORIZED_RESOURCE
	}
	return nil
}
