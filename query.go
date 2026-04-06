package fiberkit

import "github.com/gofiber/fiber/v3"

// Query parses and validates the request query before calling the typed handler.
func Query[T any](handler func(fiber.Ctx, T) error) fiber.Handler {
	return func(ctx fiber.Ctx) error {
		var query T
		if err := ctx.Bind().WithoutAutoHandling().SkipValidation(true).Query(&query); err != nil {
			return invalidInput(ctx, "query")
		}

		if err := validateInput(query); err != nil {
			return invalidValidation(ctx, err)
		}

		return handler(ctx, query)
	}
}
