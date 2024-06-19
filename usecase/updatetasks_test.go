package usecase

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"task-manager/constants"
	"task-manager/dto"
	"task-manager/model"
	"task-manager/routes"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestUpdateTasksHandler(t *testing.T) {
	testCollection := "sample-collection"
	taskID := primitive.NewObjectID()
	t.Run("UpdateTasksHandler - invalid body", func(t *testing.T) {
		invalidUpdateReqBody := ""
		res, err := json.Marshal(invalidUpdateReqBody)
		require.NoError(t, err)
		req := httptest.NewRequest(http.MethodPut, routes.UpdateTasks, bytes.NewBuffer(res))
		ctx := context.WithValue(req.Context(), constants.UserCollectionName, testCollection)
		req = req.WithContext(ctx)
		resp := httptest.NewRecorder()
		UpdateTasksHandler(resp, req)

		var response *dto.GenericResponse
		err = json.NewDecoder(resp.Body).Decode(&response)
		require.NoError(t, err)

		require.Equal(t, http.StatusBadRequest, response.Code)
		require.Equal(t, "error decoding the request body", response.Error)
		require.Equal(t, "unable to process request", response.Message)
	})

	t.Run("UpdateTasksHandler - empty data in body", func(t *testing.T) {
		emptyReqBody := &dto.TaskCreateRequestBody{
			Tasks: []*model.Task{},
		}

		res, err := json.Marshal(emptyReqBody)
		require.NoError(t, err)
		req, err := http.NewRequest(http.MethodPut, routes.UpdateTasks, bytes.NewBuffer(res))
		require.NoError(t, err)

		ctx := context.WithValue(req.Context(), constants.UserCollectionName, testCollection)
		req = req.WithContext(ctx)
		resp := httptest.NewRecorder()

		UpdateTasksHandler(resp, req)

		var response *dto.GenericResponse
		err = json.NewDecoder(resp.Body).Decode(&response)
		require.NoError(t, err)
		require.Equal(t, http.StatusBadRequest, response.Code)
		require.Equal(t, "empty body", response.Message)
	})

	t.Run("UpdateTasksHandler - success", func(t *testing.T) {
		// insert a sample task which we will update later
		BeforeSuccessTest(t, taskID, testCollection)
		reqBody := &dto.TaskCreateRequestBody{
			Tasks: []*model.Task{
				{
					ID:          taskID,
					Title:       "title-2",
					Description: "desc-2",
					Priority:    1,
					Status:      model.InProgress,
					CreatedAt:   time.Now(),
					DeadLine:    "2024-05-05",
				},
			},
		}

		res, err := json.Marshal(reqBody)
		require.NoError(t, err)
		req, err := http.NewRequest(http.MethodPut, routes.UpdateTasks, bytes.NewBuffer(res))
		require.NoError(t, err)
		ctx := context.WithValue(req.Context(), constants.UserCollectionName, testCollection)
		req = req.WithContext(ctx)

		resp := httptest.NewRecorder()
		UpdateTasksHandler(resp, req)

		var response *dto.GenericResponse
		err = json.NewDecoder(resp.Body).Decode(&response)
		require.NoError(t, err)
		require.Equal(t, http.StatusOK, response.Code)
		require.Equal(t, "request processed successfully", response.Message)
	})
}
