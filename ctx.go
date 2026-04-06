package fiberkit

import "github.com/gofiber/fiber/v2"

// Set stores a value in Fiber locals.
func Set(ctx *fiber.Ctx, key string, value any) {
	ctx.Locals(key, value)
}

// Get reads a value from Fiber locals and returns nil on missing or mismatched types.
func Get[T any](ctx *fiber.Ctx, key string) *T {
	value := ctx.Locals(key)
	if value == nil {
		return nil
	}

	typed, ok := value.(T)
	if ok {
		return &typed
	}

	pointer, ok := value.(*T)
	if ok {
		return pointer
	}

	return nil
}
