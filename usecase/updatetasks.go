package usecase

import (
	"encoding/json"
	"net/http"
	"task-manager/config"
	"task-manager/dto"
	"task-manager/repository"
	"task-manager/utils"

	"go.uber.org/zap"
)

func UpdateTasksHandler(w http.ResponseWriter, r *http.Request) {
	var requestTasks *dto.TaskCreateRequestBody

	// decode the request body
	if err := json.NewDecoder(r.Body).Decode(&requestTasks); err != nil {
		zap.L().Error("error decoding the request body", zap.Error(err))
		utils.Response(w, &dto.GenericResponse{
			Code:    http.StatusBadRequest,
			Error:   "error decoding the request body",
			Message: "unable to process request",
		}, http.StatusBadRequest)
		return
	}

	if len(*requestTasks.Tasks) == 0 {
		zap.L().Error("nothing to update")
		utils.Response(w, &dto.GenericResponse{
			Code:    http.StatusBadRequest,
			Message: "nothing to update",
		}, http.StatusBadRequest)
		return
	}

	mongoRepo := repository.MongoRepository{
		Client:     config.MongoClient,
		Collection: "sample-collection",
	}

	if err := mongoRepo.UpdateTasks(r.Context(), requestTasks); err != nil {
		zap.L().Error("error updating tasks", zap.Error(err))
		utils.Response(w, &dto.GenericResponse{
			Code:    http.StatusInternalServerError,
			Message: "unable to process request",
			Error:   err.Error(),
		}, http.StatusInternalServerError)
		return
	}
}
