package perm

import "treehole_next/models"

type Role = int

// Role enum
const (
	Admin Role = 1 << iota
	Operator
)

type PermissionOwner interface {
	GetPermission() Role
}

// GetPermission checks user permission
// Deprecated
// Example:
//
//	GetPermission(u, perm.Admin)
//	GetPermission(u, perm.Admin | perm.Operator)
func GetPermission(u PermissionOwner, p Role) bool {
	return u.GetPermission()&p != 0
}

type PermissionModel interface {
	CheckPermission(user *models.User) error
}

func CheckPermission(user *models.User, models ...PermissionModel) error {
	for _, model := range models {
		err := model.CheckPermission(user)
		if err != nil {
			return err
		}
	}
	return nil
}
