package fiberkit

import "github.com/gofiber/fiber/v2"

// Body parses and validates the request body before calling the typed handler.
func Body[T any](handler func(*fiber.Ctx, *T) error) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		var body T
		if err := ctx.BodyParser(&body); err != nil {
			return invalidInput(ctx, "body")
		}

		if err := validateInput(body); err != nil {
			return invalidValidation(ctx, err)
		}

		return handler(ctx, &body)
	}
}
