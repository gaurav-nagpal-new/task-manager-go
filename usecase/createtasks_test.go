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

func TestCreateTasksHandler(t *testing.T) {
	now := time.Now().UTC()
	testCollection := "sample-collection"
	task1ID := primitive.NewObjectID()
	task2ID := primitive.NewObjectID()
	task3ID := primitive.NewObjectID()

	deadLine := "2024-05-05"
	tasks := &dto.TaskCreateRequestBody{
		Tasks: []*model.Task{
			{
				ID:          task1ID,
				Title:       "title-1",
				Description: "desc-1",
				Priority:    1,
				Status:      model.Todo,
				CreatedAt:   now,
				DeadLine:    deadLine,
			},
			{
				ID:          task2ID,
				Title:       "title-2",
				Description: "desc-2",
				Priority:    2,
				Status:      model.Todo,
				CreatedAt:   now,
				DeadLine:    deadLine,
			},
			{
				ID:          task3ID,
				Title:       "title-3",
				Description: "desc-3",
				Priority:    3,
				Status:      model.Done,
				CreatedAt:   time.Time{},
				DeadLine:    deadLine,
			},
		},
	}
	t.Run("CreateTaskHandler - invalid body", func(t *testing.T) {
		invalidTasksBody := ""
		reqBody, err := json.Marshal(invalidTasksBody)
		require.NoError(t, err)
		req := httptest.NewRequest(http.MethodPost, routes.CreateTasks, bytes.NewBuffer(reqBody))
		ctx := context.WithValue(req.Context(), constants.UserCollectionName, testCollection)
		req = req.WithContext(ctx)
		resp := httptest.NewRecorder()

		CreateTasksHandler(resp, req)

		var response *dto.GenericResponse
		err = json.NewDecoder(resp.Body).Decode(&response)
		require.NoError(t, err)

		require.Equal(t, http.StatusBadRequest, response.Code)
		require.NotNil(t, response.Error)
		require.Equal(t, "error decoding request body", response.Message)
	})

	t.Run("CreateTaskHandler - success", func(t *testing.T) {
		reqBody, err := json.Marshal(tasks)
		require.NoError(t, err)
		req := httptest.NewRequest(http.MethodPost, routes.CreateTasks, bytes.NewBuffer(reqBody))
		ctx := context.WithValue(req.Context(), constants.UserCollectionName, testCollection)
		req = req.WithContext(ctx)
		resp := httptest.NewRecorder()

		CreateTasksHandler(resp, req)

		var response *dto.GenericResponse
		err = json.NewDecoder(resp.Body).Decode(&response)
		require.NoError(t, err)

		require.Equal(t, http.StatusOK, response.Code)
		require.Equal(t, http.StatusText(http.StatusCreated), response.Message)
	})
}
