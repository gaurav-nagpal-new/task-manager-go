package usecase

import (
	"encoding/json"
	"net/http"
	"task-manager/config"
	"task-manager/constants"
	"task-manager/dto"
	"task-manager/model"
	"task-manager/myerrors"
	"task-manager/repository"
	"task-manager/utils"
	"task-manager/validations"

	"golang.org/x/crypto/bcrypt"

	"go.uber.org/zap"
)

func SignUpHandler(w http.ResponseWriter, r *http.Request) {

	// decode the request
	var u *dto.UserCreateRequestBody
	if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
		zap.L().Error("error decoding the user create request body", zap.Error(err))
		utils.Response(w, &dto.GenericResponse{
			Code:    http.StatusBadRequest,
			Message: "unable to process request body",
			Error:   err.Error(),
		}, http.StatusBadRequest)
		return
	}

	// check if other exists with same email
	mongoRepo := repository.MongoRepository{
		Client:     config.MongoClient,
		Collection: constants.UserCollection,
	}

	// validations on user request body and send 400 in case
	if err := validations.ValidateSignUpData(u); err != nil {
		zap.L().Error("validation error", zap.Error(err))
		utils.Response(w, &dto.GenericResponse{
			Code:    http.StatusBadRequest,
			Error:   err.Error(),
			Message: "unable to process request",
		}, http.StatusBadRequest)
		return
	}
	ctx := r.Context()
	var user *model.User
	user, err := mongoRepo.FetchUserByEmail(ctx, u.User.Email)
	if err != nil && err.Error() != myerrors.NoDocumentsErr {
		zap.L().Error("error getting user data from DB", zap.Error(err))
		utils.Response(w, &dto.GenericResponse{
			Code:    http.StatusInternalServerError,
			Error:   err.Error(),
			Message: "unable to process request",
		}, http.StatusInternalServerError)
		return
	}

	// user already exists so return
	if user != nil && user.Email != "" {
		zap.L().Error("user already exists")
		utils.Response(w, &dto.GenericResponse{
			Code:    http.StatusBadRequest,
			Message: "User already exists",
		}, http.StatusBadRequest)
		return
	}

	user = &model.User{}
	// now create the user and the linked collection
	user.TaskCollection = utils.GetTaskCollectionName(u.User.Email)
	user.Name = u.User.Name
	user.Email = u.User.Email

	// hash the password before storing it
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.User.Password), bcrypt.DefaultCost)
	if err != nil {
		zap.L().Error("error hashing the password", zap.Error(err))
		utils.Response(w, &dto.GenericResponse{
			Code:    http.StatusInternalServerError,
			Message: "unable to process the request",
			Error:   err.Error(),
		}, http.StatusInternalServerError)
		return
	}

	user.Password = string(hashedPassword)

	if err := mongoRepo.CreateUser(ctx, user); err != nil {
		zap.L().Error("error creating the user", zap.Error(err))
		utils.Response(w, &dto.GenericResponse{
			Code:    http.StatusInternalServerError,
			Message: "unable to process the request",
			Error:   err.Error(),
		}, http.StatusInternalServerError)
		return
	}

	// create the collection now
	if err := mongoRepo.CreateTaskCollection(ctx, user.TaskCollection); err != nil {
		zap.L().Error("error creating the collection", zap.Error(err))
		utils.Response(w, &dto.GenericResponse{
			Code:    http.StatusInternalServerError,
			Message: "unable to process the request",
			Error:   err.Error(),
		}, http.StatusInternalServerError)
		return
	}

	// success response
	utils.Response(w, &dto.GenericResponse{
		Code:    http.StatusCreated,
		Message: "request processed successfully",
	}, http.StatusCreated)

}
