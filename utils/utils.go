package utils

import (
	"github.com/gofiber/fiber/v2"
	"github.com/opentreehole/go-common"
	"golang.org/x/exp/constraints"
	"strconv"
	"treehole_next/config"
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

func MiddlewareHasAnsweredQuestions(c *fiber.Ctx) error {
	if config.Config.Mode != "test" && config.Config.Mode != "bench" {
		return c.Next()
	}
	var user struct {
		HasAnsweredQuestions bool `json:"has_answered_questions"`
	}
	err := common.ParseJWTToken(common.GetJWTToken(c), &user)
	if err != nil {
		return err
	}
	if !user.HasAnsweredQuestions {
		return common.Forbidden("请先通过注册答题")
	}
	return c.Next()
}
