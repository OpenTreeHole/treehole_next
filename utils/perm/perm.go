package perm

type Role = int

// Role enum
const (
	Admin Role = 1 << iota
	Operator
)

type PermissionOwner interface {
	GetPermission() Role
}

// CheckPermission checks user permission
// Example:
//
//	CheckPermission(u, perm.Admin)
//	CheckPermission(u, perm.Admin | perm.Operator)
func CheckPermission(u PermissionOwner, p Role) bool {
	return u.GetPermission()&p != 0
}
