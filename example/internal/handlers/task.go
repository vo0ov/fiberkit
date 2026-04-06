package handlers

import (
	"example/internal/models"

	"github.com/gofiber/fiber/v3"
	"github.com/vo0ov/fiberkit/v3"
)

type TaskHandler struct{}

func NewTaskHandler() *TaskHandler {
	return &TaskHandler{}
}

func (h *TaskHandler) CreateTask(ctx fiber.Ctx, body models.CreateTaskRequest) error {
	return ctx.Status(fiber.StatusCreated).JSON(fiber.Map{
		"name": body.Name,
	})
}

func (h *TaskHandler) ListTasks(ctx fiber.Ctx, query models.ListTasksQuery) error {
	status := query.Status
	if status == "" {
		status = "all"
	}

	return ctx.JSON(fiber.Map{
		"status": status,
		"items": []fiber.Map{
			{"id": "task-1", "name": "Ship fiberkit"},
		},
	})
}

func (h *TaskHandler) GetTask(ctx fiber.Ctx) error {
	user := fiberkit.Get[models.User](ctx, "currentUser")
	taskID := fiberkit.Get[string](ctx, "taskID")
	if user == nil || taskID == nil {
		return ctx.SendStatus(fiber.StatusInternalServerError)
	}

	return ctx.JSON(fiber.Map{
		"id":   *taskID,
		"name": "Ship fiberkit",
		"user": user,
	})
}

func (h *TaskHandler) UpdateTask(ctx fiber.Ctx, params models.TaskParams, body models.UpdateTaskRequest) error {
	user := fiberkit.Get[models.User](ctx, "currentUser")
	if user == nil {
		return ctx.SendStatus(fiber.StatusInternalServerError)
	}

	return ctx.JSON(fiber.Map{
		"id":   params.ID,
		"name": body.Name,
		"user": user,
	})
}
