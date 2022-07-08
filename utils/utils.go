package utils

import (
	"encoding/json"
	"github.com/gofiber/fiber/v2"
)

//func InArray[T comparable](item *T, container *[]T) bool {
//	for _, i := range *container {
//		if *item == i {
//			return true
//		}
//	}
//	return false
//}

// BindJSON is a safe method to bind request body to struct
func BindJSON(c *fiber.Ctx, obj interface{}) error {
	body := c.Body()
	if len(body) == 0 {
		body, _ = json.Marshal(fiber.Map{})
	}
	return json.Unmarshal(body, obj)
}

type CanPreprocess interface {
	Preprocess() error
}

func Serialize(c *fiber.Ctx, obj CanPreprocess) error {
	err := obj.Preprocess()
	if err != nil {
		return err
	}
	return c.JSON(obj)
}

func DiffrenceSet[T comparable](mainSet []T, subSet []T) (ansSet []T) {
	tmp := map[T]bool{}
	for _, val := range subSet {
		if _, ok := tmp[val]; !ok {
			tmp[val] = true
		}
	}
	for _, val := range mainSet {
		if _, ok := tmp[val]; !ok {
			ansSet = append(ansSet, val)
		}
	}
	return ansSet
}
