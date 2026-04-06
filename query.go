package fiberkit

import "github.com/gofiber/fiber/v2"

// Query parses and validates the request query before calling the typed handler.
func Query[T any](handler func(*fiber.Ctx, *T) error) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		var query T
		if err := ctx.QueryParser(&query); err != nil {
			return invalidInput(ctx, "query")
		}

		if err := validateInput(query); err != nil {
			return invalidValidation(ctx, err)
		}

		return handler(ctx, &query)
	}
}
