package middlewares

import (
	"example/internal/models"

	"github.com/gofiber/fiber/v3"
	"github.com/vo0ov/fiberkit/v3"
)

type UserLoader struct{}

func (l *UserLoader) Load(ctx fiber.Ctx, params models.TaskParams) error {
	user := &models.User{
		ID:   "user-1",
		Role: "admin",
	}

	fiberkit.Set(ctx, "currentUser", user)
	fiberkit.Set(ctx, "taskID", params.ID)

	return ctx.Next()
}
