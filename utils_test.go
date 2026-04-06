package fiberkit

import (
	"encoding/json"
	"errors"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
)

type validationHelperPayload struct {
	Name string `validate:"required"`
}

func TestValidateInputAllowsNilValue(t *testing.T) {
	if err := validateInput(nil); err != nil {
		t.Fatalf("validateInput(nil) error = %v, want nil", err)
	}
}

func TestValidateInputAllowsNonStructValue(t *testing.T) {
	if err := validateInput("plain-string"); err != nil {
		t.Fatalf("validateInput(non-struct) error = %v, want nil", err)
	}
}

func TestValidateInputSupportsStructPointer(t *testing.T) {
	payload := &validationHelperPayload{Name: "ok"}
	if err := validateInput(payload); err != nil {
		t.Fatalf("validateInput(pointer) error = %v, want nil", err)
	}
}

func TestInvalidValidationOmitsDetailsForGenericError(t *testing.T) {
	app := fiber.New()
	app.Get("/", func(ctx *fiber.Ctx) error {
		return invalidValidation(ctx, errors.New("boom"))
	})

	req := httptest.NewRequest(fiber.MethodGet, "/", nil)
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

	if _, exists := got["details"]; exists {
		t.Fatalf("details exists = true, want false")
	}
}
