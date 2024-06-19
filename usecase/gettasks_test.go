package usecase

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"task-manager/constants"
	"task-manager/dto"
	"task-manager/routes"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGetTasksHandler(t *testing.T) {
	t.Run("GetTasksHandler - invalid status", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("%s?status=pending", routes.GetTasks), nil)
		ctx := context.WithValue(req.Context(), constants.UserCollectionName, "testCollection")
		req = req.WithContext(ctx)
		resp := httptest.NewRecorder()

		GetTasksHandler(resp, req)

		// decode the response body
		var response *dto.GenericResponse
		err := json.NewDecoder(resp.Body).Decode(&response)
		require.NoError(t, err)

		require.Equal(t, http.StatusBadRequest, response.Code)
		require.Equal(t, "unable to process request", response.Message)
		require.Equal(t, "status/priority is not valid", response.Error)
	})

	t.Run("GetTasksHandler - invalid priority", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("%s?status=Todo&priority=2", routes.GetTasks), nil)
		ctx := context.WithValue(req.Context(), constants.UserCollectionName, "testCollection")
		req = req.WithContext(ctx)
		resp := httptest.NewRecorder()

		GetTasksHandler(resp, req)

		// decode the response body
		var response *dto.GenericResponse
		err := json.NewDecoder(resp.Body).Decode(&response)
		require.NoError(t, err)

		require.Equal(t, http.StatusBadRequest, response.Code)
		require.Equal(t, "unable to process request", response.Message)
		require.Equal(t, "status/priority is not valid", response.Error)
	})

	t.Run("GetTasksHandler - invalid status", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("%s?status=Todo&priority=0", routes.GetTasks), nil)
		ctx := context.WithValue(req.Context(), constants.UserCollectionName, "testCollection")
		req = req.WithContext(ctx)
		resp := httptest.NewRecorder()

		GetTasksHandler(resp, req)

		// decode the response body
		var response *dto.GetTasksResponse
		err := json.NewDecoder(resp.Body).Decode(&response)
		require.NoError(t, err)

		require.Equal(t, http.StatusOK, response.Code)
		require.Equal(t, "Request processed successfully", response.Message)
		require.Len(t, response.Data, 0)
	})
}
