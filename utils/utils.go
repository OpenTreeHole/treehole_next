package utils

import (
	"github.com/gofiber/fiber/v2"
	"golang.org/x/exp/constraints"
	"strconv"
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

func Keys[T comparable, S any](m map[T]S) (s []T) {
	for k := range m {
		s = append(s, k)
	}
	return s
}

func Min[T constraints.Ordered](x T, y T) T {
	if x > y {
		return y
	} else {
		return x
	}
}

func StripContent(content string, contentMaxSize int) string {
	return string([]rune(content)[:Min(len([]rune(content)), contentMaxSize)])
}
