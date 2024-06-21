package dto

import "task-manager/model"

// swagger:model
type GenericResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Error   string `json:"error,omitempty"`
}

// swagger:model
type GetTasksResponse struct {
	Data    []*model.Task `json:"tasks",omitempty`
	Code    int           `json:"code"`
	Message string        `json:"message"`
	Error   string        `json:"error,omitempty"`
}
