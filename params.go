package fiberkit

import "github.com/gofiber/fiber/v2"

// Params parses and validates route params before calling the typed handler.
func Params[T any](handler func(*fiber.Ctx, *T) error) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		var params T
		if err := ctx.ParamsParser(&params); err != nil {
			return invalidInput(ctx, "params")
		}

		if err := validateInput(params); err != nil {
			return invalidValidation(ctx, err)
		}

		return handler(ctx, &params)
	}
}
