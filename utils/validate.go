package utils

import (
	"github.com/creasty/defaults"
	"github.com/go-playground/validator/v10"
	"github.com/goccy/go-json"
	"github.com/gofiber/fiber/v2"
	"reflect"
	"strings"
)

type ErrorDetailElement struct {
	Field string `json:"field"`
	Tag   string `json:"tag"`
	Value string `json:"value"`
}

type ErrorDetail []*ErrorDetailElement

func (e *ErrorDetail) Error() string {
	return "Validation Error"
}

var validate = validator.New()

func init() {
	validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]

		if name == "-" {
			return ""
		}

		return name
	})
}

func Validate(model any) error {
	errors := validate.Struct(model)
	if errors != nil {
		var errorDetail ErrorDetail
		for _, err := range errors.(validator.ValidationErrors) {
			detail := ErrorDetailElement{
				Field: err.Field(),
				Tag:   err.Tag(),
				Value: err.Param(),
			}
			errorDetail = append(errorDetail, &detail)
		}
		return &errorDetail
	}
	return nil
}

func ValidateQuery[T any](c *fiber.Ctx) (*T, error) {
	model := new(T)
	if err := c.QueryParser(model); err != nil {
		return nil, err
	}
	if err := defaults.Set(model); err != nil {
		return nil, err
	}
	return model, Validate(model)
}

// ValidateBody supports json only
func ValidateBody[T any](c *fiber.Ctx) (*T, error) {
	body := c.Body()
	model := new(T)
	if len(body) == 0 {
		return model, defaults.Set(model)
	} else {
		if err := json.Unmarshal(body, model); err != nil {
			return nil, err
		}
		if err := defaults.Set(model); err != nil {
			return nil, err
		}
		return model, Validate(model)
	}
}
