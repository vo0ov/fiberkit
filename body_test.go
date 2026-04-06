package fiberkit

import (
	"encoding/json"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gofiber/fiber/v2"
)

type createTaskRequest struct {
	Name string `json:"name" validate:"required"`
}

func TestBodyParsesValidJSON(t *testing.T) {
	app := fiber.New()
	app.Post("/tasks", Body(func(ctx *fiber.Ctx, body *createTaskRequest) error {
		return ctx.JSON(body)
	}))

	req := httptest.NewRequest(fiber.MethodPost, "/tasks", strings.NewReader(`{"name":"demo"}`))
	req.Header.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSON)

	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("app.Test() error = %v", err)
	}

	if resp.StatusCode != fiber.StatusOK {
		t.Fatalf("status = %d, want %d", resp.StatusCode, fiber.StatusOK)
	}

	var got createTaskRequest
	if err := json.NewDecoder(resp.Body).Decode(&got); err != nil {
		t.Fatalf("decode response: %v", err)
	}

	if got.Name != "demo" {
		t.Fatalf("name = %q, want %q", got.Name, "demo")
	}
}

func TestBodyRejectsInvalidJSON(t *testing.T) {
	app := fiber.New()
	app.Post("/tasks", Body(func(ctx *fiber.Ctx, body *createTaskRequest) error {
		return ctx.JSON(body)
	}))

	req := httptest.NewRequest(fiber.MethodPost, "/tasks", strings.NewReader(`{"name":`))
	req.Header.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSON)

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

	if got["error"] != "invalid body" {
		t.Fatalf("error = %v, want %q", got["error"], "invalid body")
	}
}

func TestBodyValidatesInput(t *testing.T) {
	app := fiber.New()
	app.Post("/tasks", Body(func(ctx *fiber.Ctx, body *createTaskRequest) error {
		return ctx.JSON(body)
	}))

	req := httptest.NewRequest(fiber.MethodPost, "/tasks", strings.NewReader(`{"name":""}`))
	req.Header.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSON)

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

	if details["Name"] != "required" {
		t.Fatalf("details[Name] = %v, want %q", details["Name"], "required")
	}
}
