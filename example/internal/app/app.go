package app

import (
	"example/internal/handlers"
	"example/internal/middlewares"

	"github.com/gofiber/fiber/v3"
	"github.com/vo0ov/fiberkit/v3"
)

func New() *fiber.App {
	app := fiber.New()

	// Application dependencies live in the app layer, not in fiberkit.
	handler := handlers.NewTaskHandler()
	loader := &middlewares.UserLoader{}

	// fiberkit wraps request binding and lets typed middleware and handlers stay clean.
	app.Post("/tasks", fiberkit.Body(handler.CreateTask))
	app.Get("/tasks", fiberkit.Query(handler.ListTasks))
	app.Get("/tasks/:id",
		middlewares.Auth,
		fiberkit.Params(loader.Load),
		handler.GetTask,
	)
	app.Patch("/tasks/:id",
		middlewares.Auth,
		fiberkit.Params(loader.Load),
		fiberkit.ParamsBody(handler.UpdateTask),
	)

	return app
}
