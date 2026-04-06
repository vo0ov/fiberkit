package models

type CreateTaskRequest struct {
	Name string `json:"name" validate:"required"`
}

type ListTasksQuery struct {
	Status string `query:"status"`
}

type TaskParams struct {
	ID string `params:"id" validate:"required"`
}

type UpdateTaskRequest struct {
	Name string `json:"name" validate:"required"`
}

type User struct {
	ID   string `json:"id"`
	Role string `json:"role"`
}
