package usecase

import (
	"encoding/json"
	"net/http"
	"task-manager/config"
	"task-manager/dto"
	"task-manager/jwtauth"
	"task-manager/myerrors"
	"task-manager/repository"
	"task-manager/utils"
	"task-manager/validations"

	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

func SignInHandler(w http.ResponseWriter, r *http.Request) {

	// decode the request body
	var loginData *dto.UserLoginRequestBody

	if err := json.NewDecoder(r.Body).Decode(&loginData); err != nil {
		zap.L().Error("error decoding the request body", zap.Error(err))
		utils.Response(w, &dto.GenericResponse{
			Code:    http.StatusBadRequest,
			Message: "unable to process request",
			Error:   err.Error(),
		}, http.StatusBadRequest)
		return
	}

	// validations as Email and Password is required
	if err := validations.ValidateSignInData(loginData); err != nil {
		zap.L().Error("validation failed for login", zap.Error(err))
		utils.Response(w, &dto.GenericResponse{
			Code:    http.StatusBadRequest,
			Message: "unable to process request",
			Error:   err.Error(),
		}, http.StatusBadRequest)
		return
	}

	mongoRepo := repository.MongoRepository{
		Client: config.MongoClient,
	}

	ctx := r.Context()
	// now check if the user exists in our DB
	existingUser, err := mongoRepo.FetchUserByEmail(ctx, loginData.Email)
	if err != nil {
		response := &dto.GenericResponse{
			Code:    http.StatusInternalServerError,
			Message: "unable to process request",
			Error:   err.Error(),
		}
		if err.Error() == myerrors.NoDocumentsErr {
			response.Code = http.StatusBadRequest
			response.Error = "no user found with this email"
		}

		utils.Response(w, response, response.Code)
		return
	}

	// now compare the saved password with sent password
	if err = bcrypt.CompareHashAndPassword([]byte(existingUser.Password), []byte(loginData.Password)); err != nil {
		zap.L().Error("invalid creds", zap.Error(err))
		utils.Response(w, &dto.GenericResponse{
			Code:    http.StatusUnauthorized,
			Message: "invalid creds",
		}, http.StatusUnauthorized)
		return
	}

	// Create JWT token here and sent it to the headers (need to use these tokens in all endpoints(middleware))
	token, err := jwtauth.GenerateJWTToken(loginData.Email)
	if err != nil {
		zap.L().Error("error creating jwt token", zap.Error(err))
		utils.Response(w, &dto.GenericResponse{
			Code:    http.StatusInternalServerError,
			Message: "unable to process request",
			Error:   err.Error(),
		}, http.StatusInternalServerError)
		return
	}

	// send token in headers
	w.Header().Set("Authorization", "Bearer "+token)

	// otherwise success
	utils.Response(w, &dto.GenericResponse{
		Code:    http.StatusOK,
		Message: "login successfull",
	}, http.StatusOK)
}
