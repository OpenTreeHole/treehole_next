// Package models contains database models
package models

type Map = map[string]interface{}

type Models interface {
	Division | Hole | Floor | Tag | User | Report |
		[]Division | []Hole | []Floor | []Tag | []User | []Report
}
