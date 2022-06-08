package models

type ResponseModel interface {
	DivisionResponse | Hole | Floor | Tag | User
}

type ResponseModelSlice[T ResponseModel] []T

type ResponseModelFull interface {
	ResponseModel
}

type ModelsSingle interface {
	Division | Hole | Floor | Tag | User
}

type Models interface {
	ModelsSingle
}
