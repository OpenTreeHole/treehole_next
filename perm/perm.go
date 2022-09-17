package perm

// PermissionType enum
const (
	Admin = 1 << iota
	Operator
)

type PermissionType = int
