package usecase

import (
	"encoding/json"
	"net/http"
	"task-manager/config"
	"task-manager/constants"
	"task-manager/dto"
	"task-manager/repository"
	"task-manager/utils"

	"go.uber.org/zap"
)

func UpdateTasksHandler(w http.ResponseWriter, r *http.Request) {
	// swagger:operation PUT /tasks/update Tasks Update Tasks
	// Update Tasks
	//
	// Update Tasks
	// ---
	// tags:
	// - Tasks
	// produces:
	// - application/json
	// parameters:
	// - name: tasks
	//   in: body
	//   required: true
	//   schema:
	//     $ref: '#definitions/TaskCreateRequestBody'
	// responses:
	//   '200':
	//     schema:
	//       $ref: '#/definitions/GenericResponse'
	//   '400':
	//     schema:
	//       type: object
	//       properties:
	//         code:
	//           type: int
	//           example: 400
	//         error:
	//           type: string
	//           example: "error decoding request body"
	//         message:
	//           type: string
	//           example: "unable to process request"
	//   '500':
	//     schema:
	//       type: object
	//       properties:
	//         code:
	//           type: int
	//           example: 500
	//         error:
	//           type: string
	//           example: "Unable to process the request"
	//         message:
	//           type: string
	//           example: "An error occurred"
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

	if len(requestTasks.Tasks) == 0 {
		zap.L().Error("empty body")
		utils.Response(w, &dto.GenericResponse{
			Code:    http.StatusBadRequest,
			Message: "empty body",
		}, http.StatusBadRequest)
		return
	}

	mongoRepo := repository.MongoRepository{
		Client:     config.MongoClient,
		Collection: r.Context().Value(constants.UserCollectionName).(string),
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

	utils.Response(w, &dto.GenericResponse{
		Code:    http.StatusOK,
		Message: "request processed successfully",
	}, http.StatusOK)
}
