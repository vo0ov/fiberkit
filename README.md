# FiberKit ✨

Lightweight typed request binding helpers for Fiber.

`fiberkit` makes Fiber routes easier to read by parsing request data, validating it, and passing typed arguments into handlers and typed middleware.

## Before ❌

```go
app.Post("/tasks", func(ctx fiber.Ctx) error {
	var body CreateTaskRequest
	if err := ctx.Bind().Body(&body); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid body",
		})
	}

	if err := validator.New().Struct(body); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "validation failed",
		})
	}

	return ctx.JSON(body)
})
```

## After ✅

```go
app.Post("/tasks", fiberkit.Body(CreateTask))

func CreateTask(ctx fiber.Ctx, body CreateTaskRequest) error {
	return ctx.JSON(body)
}
```

## Install 🚀

```bash
go get github.com/vo0ov/fiberkit/v3
```

## API 🛠️

```go
func Body[T any](handler func(fiber.Ctx, T) error) fiber.Handler
func Query[T any](handler func(fiber.Ctx, T) error) fiber.Handler
func Params[T any](handler func(fiber.Ctx, T) error) fiber.Handler
func ParamsBody[P any, B any](handler func(fiber.Ctx, P, B) error) fiber.Handler

func Set(ctx fiber.Ctx, key string, value any)
func Get[T any](ctx fiber.Ctx, key string) *T
```

## Usage 💡

### Body

```go
type CreateTaskRequest struct {
	Name string `json:"name" validate:"required"`
}

app.Post("/tasks", fiberkit.Body(CreateTask))

func CreateTask(ctx fiber.Ctx, body CreateTaskRequest) error {
	return ctx.Status(fiber.StatusCreated).JSON(body)
}
```

### Query

```go
type ListTasksQuery struct {
	Status string `query:"status"`
}

app.Get("/tasks", fiberkit.Query(ListTasks))

func ListTasks(ctx fiber.Ctx, query ListTasksQuery) error {
	return ctx.JSON(query)
}
```

### Params

```go
type GetTaskParams struct {
	ID string `uri:"id" validate:"required"`
}

app.Get("/tasks/:id", fiberkit.Params(GetTask))

func GetTask(ctx fiber.Ctx, params GetTaskParams) error {
	return ctx.JSON(params)
}
```

Note: on Fiber v3 route params are bound with `uri:"..."` tags.

### Params + Body

```go
type UpdateTaskParams struct {
	ID string `uri:"id" validate:"required"`
}

type UpdateTaskRequest struct {
	Name string `json:"name" validate:"required"`
}

app.Patch("/tasks/:id", fiberkit.ParamsBody(UpdateTask))

func UpdateTask(ctx fiber.Ctx, params UpdateTaskParams, body UpdateTaskRequest) error {
	return ctx.JSON(fiber.Map{
		"id":   params.ID,
		"name": body.Name,
	})
}
```

## Typed middleware 🤝

The same wrappers also work for middleware.

```go
type User struct {
	ID string
}

type GetUserParams struct {
	ID string `uri:"id" validate:"required"`
}

type UserLoader struct{}

func (l *UserLoader) Load(ctx fiber.Ctx, params GetUserParams) error {
	fiberkit.Set(ctx, "targetUser", &User{ID: params.ID})
	return ctx.Next()
}
```

Usage:

```go
app.Get("/users/:id",
	authMiddleware,
	fiberkit.Params(userLoader.Load),
	handler.GetUser,
)
```

## Context helpers 🧠

Use `Set` and `Get` to pass data through `c.Locals` without boilerplate.

```go
fiberkit.Set(ctx, "currentUser", &User{ID: "42"})

user := fiberkit.Get[User](ctx, "currentUser")
if user == nil {
	return ctx.SendStatus(fiber.StatusUnauthorized)
}

return ctx.JSON(user)
```

`Get[T]` returns `nil` if the key is missing or the type does not match.

## Validation ✅

Validation is built in via `go-playground/validator/v10`.

If your struct has `validate` tags, wrappers validate it automatically.

```go
type CreateTaskRequest struct {
	Name string `json:"name" validate:"required"`
}
```

Validation error response:

```json
{
    "error": "validation failed",
    "details": {
        "Name": "required"
    }
}
```

## Errors ⚠️

Parse errors return `400 Bad Request`:

```json
{
    "error": "invalid body"
}
```

```json
{
    "error": "invalid query"
}
```

```json
{
    "error": "invalid params"
}
```

Validation errors also return `400 Bad Request`.

## Notes 📌

- `fiberkit` is not a framework
- it does not do dependency injection
- it does not know about auth, services, or repositories
- it only handles parsing, validation, and typed request binding

## Contributing 🤝

Contributions are welcome.

Before opening a pull request:

- check whether the idea fits the library scope
- keep the API small and explicit
- add or update tests for behavior changes
- avoid framework-style abstractions and hidden magic

If you want to discuss a change first, open an issue and describe the problem you are trying to solve.

## License 📄

FiberKit is licensed under the [Apache-2.0](./LICENSE) License. This means you can use, modify, and distribute it freely, but it comes with no warranty. See the LICENSE file for details.
