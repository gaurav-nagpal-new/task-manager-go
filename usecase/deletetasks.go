package usecase

import (
	"net/http"
	"task-manager/config"
	"task-manager/constants"
	"task-manager/dto"
	"task-manager/repository"
	"task-manager/utils"

	"go.uber.org/zap"
)

func DeleteTasksHandler(w http.ResponseWriter, r *http.Request) {
	// get the ID from query parameter from query string

	mongoRepo := repository.MongoRepository{
		Client:     config.MongoClient,
		Collection: r.Context().Value(constants.UserCollectionName).(string),
	}

	id := r.URL.Query().Get(constants.ID)
	if id == "" {
		zap.L().Error("id not passed in query params")
		utils.Response(w, &dto.GenericResponse{
			Code:    http.StatusBadRequest,
			Message: "id field is required in query string",
			Error:   "unable to process request",
		}, http.StatusBadRequest)
		return
	}

	if err := mongoRepo.DeleteTaskFromID(r.Context(), id); err != nil {
		zap.L().Error("error deleting task from DB", zap.Error(err))
		utils.Response(w, &dto.GenericResponse{
			Code:    http.StatusInternalServerError,
			Message: "internal server error",
			Error:   err.Error(),
		}, http.StatusInternalServerError)
		return
	}

	utils.Response(w, &dto.GenericResponse{
		Code:    http.StatusOK,
		Message: "processed request successfully",
	}, http.StatusOK)
}
