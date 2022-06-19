package models

type Models interface {
	Division | Hole | Floor | Tag | User |
		[]Division | []Hole | []Floor | []Tag | []User
}

type Model interface {
	Division | Hole | Floor | Tag | User
}
