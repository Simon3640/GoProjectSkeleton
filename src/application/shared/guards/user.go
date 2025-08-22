package guards

import (
	"errors"
	usecase "gormgoskeleton/src/application/shared/use_case"
	"gormgoskeleton/src/domain/models"
)

func RoleGuard(allowedRoles ...string) usecase.Guard {
	return func(user models.UserWithRole, input any) error {
		for _, role := range allowedRoles {
			if user.GetRoleKey() == role {
				return nil
			}
		}
		return errors.New("user is not allowed to perform this action")
	}
}

// Partial Resource with UserID
func UserResourceGuard[T models.HasUserID]() usecase.Guard {
	return func(user models.UserWithRole, input any) error {
		resource, ok := input.(T)
		if !ok {
			return errors.New("invalid resource type")
		}
		if user.ID != resource.GetUserID() {
			return errors.New("user is not allowed to access this resource")
		}
		return nil
	}
}

func UserGetItSelf(user models.UserWithRole, input any) error {
	id, ok := input.(uint)
	if !ok {
		return errors.New("invalid input type")
	}
	if user.ID != id {
		return errors.New("user is not allowed to access this resource")
	}
	return nil
}
