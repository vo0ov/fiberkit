package fiberkit

import (
	"encoding/json"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gofiber/fiber/v3"
)

type updateTaskParams struct {
	ID string `uri:"id" validate:"required"`
}

type numericUpdateTaskParams struct {
	ID int `uri:"id"`
}

type validatedUpdateTaskParams struct {
	ID string `uri:"id" validate:"min=5"`
}

type updateTaskRequest struct {
	Name string `json:"name" validate:"required"`
}

func TestParamsBodyParsesParamsAndBody(t *testing.T) {
	app := fiber.New()
	app.Patch("/tasks/:id", ParamsBody(func(ctx fiber.Ctx, params updateTaskParams, body updateTaskRequest) error {
		return ctx.JSON(fiber.Map{
			"id":   params.ID,
			"name": body.Name,
		})
	}))

	req := httptest.NewRequest(fiber.MethodPatch, "/tasks/99", strings.NewReader(`{"name":"updated"}`))
	req.Header.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSON)

	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("app.Test() error = %v", err)
	}

	if resp.StatusCode != fiber.StatusOK {
		t.Fatalf("status = %d, want %d", resp.StatusCode, fiber.StatusOK)
	}

	var got map[string]any
	if err := json.NewDecoder(resp.Body).Decode(&got); err != nil {
		t.Fatalf("decode response: %v", err)
	}

	if got["id"] != "99" {
		t.Fatalf("id = %v, want %q", got["id"], "99")
	}

	if got["name"] != "updated" {
		t.Fatalf("name = %v, want %q", got["name"], "updated")
	}
}

func TestParamsBodyValidatesBody(t *testing.T) {
	app := fiber.New()
	app.Patch("/tasks/:id", ParamsBody(func(ctx fiber.Ctx, params updateTaskParams, body updateTaskRequest) error {
		return ctx.JSON(body)
	}))

	req := httptest.NewRequest(fiber.MethodPatch, "/tasks/99", strings.NewReader(`{"name":""}`))
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
}

func TestParamsBodyRejectsInvalidParams(t *testing.T) {
	app := fiber.New()
	app.Patch("/tasks/:id", ParamsBody(func(ctx fiber.Ctx, params numericUpdateTaskParams, body updateTaskRequest) error {
		return ctx.JSON(body)
	}))

	req := httptest.NewRequest(fiber.MethodPatch, "/tasks/not-a-number", strings.NewReader(`{"name":"updated"}`))
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

	if got["error"] != "invalid params" {
		t.Fatalf("error = %v, want %q", got["error"], "invalid params")
	}
}

func TestParamsBodyRejectsInvalidJSON(t *testing.T) {
	app := fiber.New()
	app.Patch("/tasks/:id", ParamsBody(func(ctx fiber.Ctx, params updateTaskParams, body updateTaskRequest) error {
		return ctx.JSON(body)
	}))

	req := httptest.NewRequest(fiber.MethodPatch, "/tasks/99", strings.NewReader(`{"name":`))
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

func TestParamsBodyValidatesParams(t *testing.T) {
	app := fiber.New()
	app.Patch("/tasks/:id", ParamsBody(func(ctx fiber.Ctx, params validatedUpdateTaskParams, body updateTaskRequest) error {
		return ctx.JSON(body)
	}))

	req := httptest.NewRequest(fiber.MethodPatch, "/tasks/abc", strings.NewReader(`{"name":"updated"}`))
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

	if details["ID"] != "min" {
		t.Fatalf("details[ID] = %v, want %q", details["ID"], "min")
	}
}
