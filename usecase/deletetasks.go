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
	// swagger:operation DELETE /tasks/delete Tasks Delete Tasks
	// Delete Tasks
	//
	// Delete Tasks
	// ---
	// tags:
	// - Tasks
	// produces:
	// - application/json
	// parameters:
	// - name: id
	//   in: query
	//   required: true
	//   schema:
	//	 	type: string
	//		example: Z34b
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
