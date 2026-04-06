package fiberkit

import "github.com/gofiber/fiber/v2"

// ParamsBody parses and validates both route params and request body.
func ParamsBody[P any, B any](handler func(*fiber.Ctx, *P, *B) error) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		var params P
		if err := ctx.ParamsParser(&params); err != nil {
			return invalidInput(ctx, "params")
		}

		var body B
		if err := ctx.BodyParser(&body); err != nil {
			return invalidInput(ctx, "body")
		}

		if err := validateInput(params); err != nil {
			return invalidValidation(ctx, err)
		}

		if err := validateInput(body); err != nil {
			return invalidValidation(ctx, err)
		}

		return handler(ctx, &params, &body)
	}
}
