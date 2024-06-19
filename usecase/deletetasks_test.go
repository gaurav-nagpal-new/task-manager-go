package usecase

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"task-manager/config"
	"task-manager/constants"
	"task-manager/dto"
	"task-manager/model"
	"task-manager/repository"
	"task-manager/routes"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestDeleteTasksHandler(t *testing.T) {
	taskID := primitive.NewObjectID()
	testCollection := "sample-collection"
	t.Run("DeleteTasksHandler - query parameter missing", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodDelete, routes.DeleteTasks, nil)
		ctx := context.WithValue(req.Context(), constants.UserCollectionName, testCollection)
		req = req.WithContext(ctx)
		resp := httptest.NewRecorder()

		DeleteTasksHandler(resp, req)

		var response *dto.GenericResponse
		err := json.NewDecoder(resp.Body).Decode(&response)
		require.NoError(t, err)

		require.Equal(t, http.StatusBadRequest, response.Code)
		require.Equal(t, "id field is required in query string", response.Message)
		require.Equal(t, "unable to process request", response.Error)
	})

	t.Run("DeleteTaskHandler - no document deleted", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodDelete, fmt.Sprintf("%s?id=%s", routes.DeleteTasks, taskID.Hex()), nil)

		ctx := context.WithValue(req.Context(), constants.UserCollectionName, testCollection)
		req = req.WithContext(ctx)

		resp := httptest.NewRecorder()

		DeleteTasksHandler(resp, req)

		var response *dto.GenericResponse
		err := json.NewDecoder(resp.Body).Decode(&response)
		require.NoError(t, err)

		require.Equal(t, http.StatusInternalServerError, response.Code)
		require.Equal(t, "internal server error", response.Message)
		require.Equal(t, "no document found to delete with the id passed", response.Error)
	})

	t.Run("DeleteTaskHandler - success", func(t *testing.T) {
		BeforeSuccessTest(t, taskID, testCollection)

		req := httptest.NewRequest(http.MethodDelete, fmt.Sprintf("%s?id=%s", routes.DeleteTasks, taskID.Hex()), nil)
		ctx := context.WithValue(req.Context(), constants.UserCollectionName, testCollection)
		req = req.WithContext(ctx)
		resp := httptest.NewRecorder()

		DeleteTasksHandler(resp, req)

		var response *dto.GenericResponse
		err := json.NewDecoder(resp.Body).Decode(&response)
		require.NoError(t, err)

		require.Equal(t, http.StatusOK, response.Code)
		require.Equal(t, "processed request successfully", response.Message)
	})
}

func BeforeSuccessTest(t *testing.T, taskID primitive.ObjectID, testCollection string) {
	// insert one document which will be deleted
	mongoRepo := repository.MongoRepository{
		Client:     config.MongoClient,
		Collection: testCollection,
	}

	sampleTasks := []*model.Task{
		{
			ID:          taskID,
			Title:       "title-1",
			Description: "desc-1",
			Priority:    1,
			Status:      model.Todo,
			CreatedAt:   time.Now(),
			DeadLine:    "2024-05-05",
		},
	}

	err := mongoRepo.CreateTasks(context.Background(), sampleTasks)
	require.NoError(t, err)
}
