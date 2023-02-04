package utils

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
)

type CanPreprocess interface {
	Preprocess(c *fiber.Ctx) error
}

func Serialize(c *fiber.Ctx, obj CanPreprocess) error {
	err := obj.Preprocess(c)
	if err != nil {
		return err
	}
	return c.JSON(obj)
}

func RegText2IntArray(IDs [][]string) ([]int, error) {
	ansIDs := make([]int, 0)
	for _, v := range IDs {
		id, err := strconv.Atoi(v[1])
		if err != nil {
			return nil, err
		}
		ansIDs = append(ansIDs, id)
	}
	return ansIDs, nil
}

type PointerConstraint[T any] interface {
	*T
}

func ValueCopy[T any, PT PointerConstraint[T]](value PT) PT {
	newValue := new(T)
	*newValue = *value
	return newValue
}

func Keys[T comparable, S any](m map[T]S) (s []T) {
	for k := range m {
		s = append(s, k)
	}
	return s
}

type numbers interface {
	int | uint | int8 | uint8 |
		int16 | uint16 | int32 | uint32 |
		int64 | uint64 | float32 | float64
}

func Min[T numbers](x T, y T) T {
	if x > y {
		return y
	} else {
		return x
	}
}

func StripContent(content string, contentMaxSize int) string {
	return content[:Min(len(content), contentMaxSize)]
}
