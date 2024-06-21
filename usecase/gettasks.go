package usecase

import (
	"net/http"
	"task-manager/config"
	"task-manager/constants"
	"task-manager/dto"
	"task-manager/model"
	"task-manager/repository"
	"task-manager/utils"

	"go.uber.org/zap"
)

func GetTasksHandler(w http.ResponseWriter, r *http.Request) {
	// swagger:operation GET /tasks/get Tasks GetTasks
	// Get Tasks
	//
	// Get Tasks
	// ---
	// tags:
	// - Tasks
	// produces:
	// - application/json
	// parameters:
	// - name: status
	//   in: query
	//   required: false
	//   schema:
	//     type: string
	//     example: "InProgress"
	// - name: priority
	//   in: query
	//   required: false
	//   schema:
	//     type: int
	//     example: 1
	// responses:
	//   '200':
	//     schema:
	//       $ref: '#/definitions/GetTasksResponse'
	//   '400':
	//     schema:
	//       type: object
	//       properties:
	//         code:
	//           type: int
	//           example: 400
	//         error:
	//           type: string
	//           example: "status/priority sent in query string is not valid"
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

	// Need to have some query params - created_at, priority and status
	mongoRepo := repository.MongoRepository{
		Client:     config.MongoClient,
		Collection: r.Context().Value(constants.UserCollectionName).(string),
	}

	status := r.URL.Query().Get("status")
	priority := utils.ConvertToInt(r.URL.Query().Get("priority"))

	// filter based on priority, status
	if status != "" || priority != 0 {
		if !model.IsValidStatus(r.URL.Query().Get("status")) || !model.IsValidPriority(priority) {
			zap.L().Error("status/priority sent in query string is not valid")
			utils.Response(w, &dto.GenericResponse{
				Code:    http.StatusBadRequest,
				Message: "unable to process request",
				Error:   "status/priority is not valid",
			}, http.StatusBadRequest)
			return
		}
	}

	result, err := mongoRepo.FetchTasks(r.Context(), status, priority)
	if err != nil {
		utils.Response(w, &dto.GenericResponse{
			Code:    http.StatusInternalServerError,
			Message: "Unable to process request",
			Error:   err.Error(),
		}, http.StatusInternalServerError)
		return
	}

	utils.Response(w, &dto.GetTasksResponse{
		Data:    result,
		Message: "Request processed successfully",
		Code:    http.StatusOK,
	}, http.StatusOK)
}
