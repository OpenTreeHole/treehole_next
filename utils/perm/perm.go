package perm

type Permission = int

// Permission enum
const (
	Admin Permission = 1 << iota
	Operator
)

type PermissionOwner interface {
	GetPermission() Permission
}

// CheckPermission checks user permission
// Example:
//  CheckPermission(u, perm.Admin)
//  CheckPermission(u, perm.Admin | perm.Operator)
//
func CheckPermission(u PermissionOwner, p Permission) bool {
	return u.GetPermission()&p != 0
}
