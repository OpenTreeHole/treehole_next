package schemas

import "treehole_next/models"

type ResponseModels interface {
	DivisionResponse | []DivisionResponse |
		models.Hole | []models.Hole |
		models.Floor | []models.Floor |
		models.Tag | []models.Tag |
		models.User | []models.User
}
