package usecase

import (
	"encoding/json"
	"net/http"
	"task-manager/config"
	"task-manager/constants"
	"task-manager/dto"
	"task-manager/repository"
	"time"

	"task-manager/utils"

	"go.uber.org/zap"
)

func CreateTasksHandler(w http.ResponseWriter, r *http.Request) {
	tasks := &dto.TaskCreateRequestBody{}
	if err := json.NewDecoder(r.Body).Decode(&tasks); err != nil {
		zap.L().Error("error decoding request body", zap.Error(err))
		response := &dto.GenericResponse{
			Code:    http.StatusBadRequest,
			Message: "error decoding request body",
			Error:   err.Error(),
		}
		utils.Response(w, response, http.StatusBadRequest)
		return
	}

	// call mongoDB function here to create tasks in DB and send response
	mongoRepo := repository.MongoRepository{
		Client:     config.MongoClient,
		Collection: r.Context().Value(constants.UserCollectionName).(string),
	}

	// set createdAt
	tasksData := tasks.Tasks
	now := time.Now()
	for i := 0; i < len(tasksData); i++ {
		tasksData[i].CreatedAt = now
	}

	if err := mongoRepo.CreateTasks(r.Context(), tasks.Tasks); err != nil {
		zap.L().Error("error creating tasks in db", zap.Error(err))
		utils.Response(w, &dto.GenericResponse{
			Code:    http.StatusInternalServerError,
			Message: http.StatusText(http.StatusInternalServerError),
		}, http.StatusInternalServerError)
		return
	}

	utils.Response(w, &dto.GenericResponse{
		Code:    http.StatusOK,
		Message: http.StatusText(http.StatusCreated),
	}, http.StatusCreated)
}
