package dto

import "task-manager/model"

type TaskCreateRequestBody struct {
	Tasks *[]model.Task `json:"data"`
}
