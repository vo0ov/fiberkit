package fiberkit

import (
	"encoding/json"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v3"
)

type getUserParams struct {
	ID string `uri:"id" validate:"required"`
}

type numericUserParams struct {
	ID int `uri:"id"`
}

type validatedUserParams struct {
	ID string `uri:"id" validate:"min=5"`
}

func TestParamsWorkForMiddleware(t *testing.T) {
	app := fiber.New()
	app.Get("/users/:id",
		Params(func(ctx fiber.Ctx, params getUserParams) error {
			Set(ctx, "targetUserID", params.ID)
			return ctx.Next()
		}),
		func(ctx fiber.Ctx) error {
			id := Get[string](ctx, "targetUserID")
			if id == nil {
				return ctx.SendStatus(fiber.StatusInternalServerError)
			}
			return ctx.SendString(*id)
		},
	)

	req := httptest.NewRequest(fiber.MethodGet, "/users/42", nil)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("app.Test() error = %v", err)
	}

	if resp.StatusCode != fiber.StatusOK {
		t.Fatalf("status = %d, want %d", resp.StatusCode, fiber.StatusOK)
	}

	body := make([]byte, 2)
	if _, err := resp.Body.Read(body); err != nil && err.Error() != "EOF" {
		t.Fatalf("read body: %v", err)
	}

	if string(body) != "42" {
		t.Fatalf("body = %q, want %q", string(body), "42")
	}
}

func TestParamsRejectsInvalidParamType(t *testing.T) {
	app := fiber.New()
	app.Get("/users/:id", Params(func(ctx fiber.Ctx, params numericUserParams) error {
		return ctx.SendStatus(fiber.StatusNoContent)
	}))

	req := httptest.NewRequest(fiber.MethodGet, "/users/not-a-number", nil)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("app.Test() error = %v", err)
	}

	if resp.StatusCode != fiber.StatusBadRequest {
		t.Fatalf("status = %d, want %d", resp.StatusCode, fiber.StatusBadRequest)
	}

	var got map[string]any
	if err := json.NewDecoder(resp.Body).Decode(&got); err != nil {
		t.Fatalf("decode response: %v", err)
	}

	if got["error"] != "invalid params" {
		t.Fatalf("error = %v, want %q", got["error"], "invalid params")
	}
}

func TestParamsValidatesInput(t *testing.T) {
	app := fiber.New()
	app.Get("/users/:id", Params(func(ctx fiber.Ctx, params validatedUserParams) error {
		return ctx.SendStatus(fiber.StatusNoContent)
	}))

	req := httptest.NewRequest(fiber.MethodGet, "/users/abc", nil)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("app.Test() error = %v", err)
	}

	if resp.StatusCode != fiber.StatusBadRequest {
		t.Fatalf("status = %d, want %d", resp.StatusCode, fiber.StatusBadRequest)
	}

	var got map[string]any
	if err := json.NewDecoder(resp.Body).Decode(&got); err != nil {
		t.Fatalf("decode response: %v", err)
	}

	if got["error"] != "validation failed" {
		t.Fatalf("error = %v, want %q", got["error"], "validation failed")
	}

	details, ok := got["details"].(map[string]any)
	if !ok {
		t.Fatalf("details type = %T, want map[string]any", got["details"])
	}

	if details["ID"] != "min" {
		t.Fatalf("details[ID] = %v, want %q", details["ID"], "min")
	}
}
