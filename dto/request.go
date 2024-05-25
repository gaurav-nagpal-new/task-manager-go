package dto

import "task-manager/model"

type TaskCreateRequestBody struct {
	Tasks *[]model.Task `json:"data"`
}

type UserCreateRequestBody struct {
	User *model.User `json:"user"`
}

type UserLoginRequestBody struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
