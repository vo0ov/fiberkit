package fiberkit

import (
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
)

type testUser struct {
	ID string
}

func TestContextSetGet(t *testing.T) {
	app := fiber.New()
	app.Get("/", func(ctx *fiber.Ctx) error {
		Set(ctx, "user", &testUser{ID: "7"})

		user := Get[testUser](ctx, "user")
		if user == nil {
			return ctx.SendStatus(fiber.StatusInternalServerError)
		}

		if user.ID != "7" {
			return ctx.SendStatus(fiber.StatusConflict)
		}

		if Get[int](ctx, "user") != nil {
			return ctx.SendStatus(fiber.StatusTeapot)
		}

		if Get[string](ctx, "missing") != nil {
			return ctx.SendStatus(fiber.StatusBadRequest)
		}

		return ctx.SendStatus(fiber.StatusNoContent)
	})

	req := httptest.NewRequest(fiber.MethodGet, "/", nil)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("app.Test() error = %v", err)
	}

	if resp.StatusCode != fiber.StatusNoContent {
		t.Fatalf("status = %d, want %d", resp.StatusCode, fiber.StatusNoContent)
	}
}
