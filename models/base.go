// Package models contains database models
package models

type Map = map[string]interface{}

type Models interface {
	Division | Hole | Floor | Tag | User | Report | Message |
		Divisions | Holes | Floors | Tags | Users | Reports | Messages
}

type MessageModel struct {
	Message string `json:"message"`
}
