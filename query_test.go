package fiberkit

import (
	"encoding/json"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v3"
)

type listTasksQuery struct {
	Page int `query:"page" validate:"required"`
}

func TestQueryParsesValues(t *testing.T) {
	app := fiber.New()
	app.Get("/tasks", Query(func(ctx fiber.Ctx, query listTasksQuery) error {
		return ctx.JSON(query)
	}))

	req := httptest.NewRequest(fiber.MethodGet, "/tasks?page=3", nil)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("app.Test() error = %v", err)
	}

	if resp.StatusCode != fiber.StatusOK {
		t.Fatalf("status = %d, want %d", resp.StatusCode, fiber.StatusOK)
	}

	var got listTasksQuery
	if err := json.NewDecoder(resp.Body).Decode(&got); err != nil {
		t.Fatalf("decode response: %v", err)
	}

	if got.Page != 3 {
		t.Fatalf("page = %d, want %d", got.Page, 3)
	}
}

func TestQueryValidatesInput(t *testing.T) {
	app := fiber.New()
	app.Get("/tasks", Query(func(ctx fiber.Ctx, query listTasksQuery) error {
		return ctx.JSON(query)
	}))

	req := httptest.NewRequest(fiber.MethodGet, "/tasks", nil)
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

	if details["Page"] != "required" {
		t.Fatalf("details[Page] = %v, want %q", details["Page"], "required")
	}
}

func TestQueryRejectsInvalidQueryInput(t *testing.T) {
	app := fiber.New()
	app.Get("/tasks", Query(func(ctx fiber.Ctx, query listTasksQuery) error {
		return ctx.JSON(query)
	}))

	req := httptest.NewRequest(fiber.MethodGet, "/tasks?page=invalid", nil)
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

	if got["error"] != "invalid query" {
		t.Fatalf("error = %v, want %q", got["error"], "invalid query")
	}
}
