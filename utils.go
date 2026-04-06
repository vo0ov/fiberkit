package fiberkit

import (
	"errors"
	"reflect"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

func validateInput(value any) error {
	typ := reflect.TypeOf(value)
	if typ == nil {
		return nil
	}

	if typ.Kind() == reflect.Pointer {
		typ = typ.Elem()
	}

	if typ.Kind() != reflect.Struct {
		return nil
	}

	return validator.New().Struct(value)
}

func invalidInput(ctx *fiber.Ctx, target string) error {
	return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
		"error": "invalid " + target,
	})
}

func invalidValidation(ctx *fiber.Ctx, err error) error {
	payload := fiber.Map{
		"error": "validation failed",
	}

	if validationErrs, ok := errors.AsType[validator.ValidationErrors](err); ok {
		details := fiber.Map{}
		for _, fieldErr := range validationErrs {
			details[fieldErr.Field()] = fieldErr.Tag()
		}
		if len(details) > 0 {
			payload["details"] = details
		}
	}

	return ctx.Status(fiber.StatusBadRequest).JSON(payload)
}
