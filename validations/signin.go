package validations

import (
	"fmt"
	"task-manager/dto"
	"task-manager/myerrors"
)

func ValidateSignInData(loginData *dto.UserLoginRequestBody) error {
	if loginData.Email == "" {
		return fmt.Errorf(myerrors.RequireFieldErrorString, "email")
	}

	if loginData.Password == "" {
		return fmt.Errorf(myerrors.RequireFieldErrorString, "password")
	}
	return nil
}

func ValidateSignUpData(signUpData *dto.UserCreateRequestBody) error {
	if signUpData.User.Name == "" {
		return fmt.Errorf(myerrors.RequireFieldErrorString, "name")
	}
	if signUpData.User.Email == "" {
		return fmt.Errorf(myerrors.RequireFieldErrorString, "email")
	}
	if signUpData.User.Password == "" {
		return fmt.Errorf(myerrors.RequireFieldErrorString, "password")

	}
	return nil
}
