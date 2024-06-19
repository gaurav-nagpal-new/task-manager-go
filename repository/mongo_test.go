package repository

import (
	"context"
	"task-manager/config"
	"task-manager/dto"
	"task-manager/model"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestMongoFunctions(t *testing.T) {
	ctx := context.Background()
	now := time.Now().UTC()
	testCollection := "sample-collection"
	task1ID := primitive.NewObjectID()
	task2ID := primitive.NewObjectID()
	task3ID := primitive.NewObjectID()

	deadLine := "2024-05-05"
	mongoRepo := MongoRepository{
		Client:     config.MongoClient,
		Collection: testCollection,
	}

	t.Run("Create Task Collection", func(t *testing.T) {
		err := mongoRepo.CreateTaskCollection(ctx, testCollection)
		require.NoError(t, err)
	})

	t.Run("Test Task Functions", func(t *testing.T) {
		t.Run("CreateTasks", func(t *testing.T) {
			sampleTasks := []*model.Task{
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
			}

			err := mongoRepo.CreateTasks(ctx, sampleTasks)
			require.NoError(t, err)
		})

		t.Run("FetchTasks - without filter", func(t *testing.T) {
			result, err := mongoRepo.FetchTasks(ctx, "", 0)
			require.NoError(t, err)
			require.Len(t, result, 3)
		})

		t.Run("FetchTasks - with status filter", func(t *testing.T) {
			result, err := mongoRepo.FetchTasks(ctx, model.Done, 0)
			require.NoError(t, err)
			require.Len(t, result, 1)
			require.Equal(t, model.Done, result[0].Status)
		})

		t.Run("UpdateTasks", func(t *testing.T) {
			updateReqBody := &dto.TaskCreateRequestBody{
				Tasks: []*model.Task{
					{
						ID:          task1ID,
						Title:       "title-1",
						Description: "desc-1",
						Priority:    1,
						Status:      model.Done,
						CreatedAt:   now,
						DeadLine:    deadLine,
					},
					{
						ID:          task2ID,
						Title:       "title-2",
						Description: "desc-2",
						Priority:    2,
						Status:      model.Done,
						CreatedAt:   now,
						DeadLine:    deadLine,
					},
					{
						ID:          task3ID,
						Title:       "title-3",
						Description: "desc-3",
						Priority:    3,
						Status:      model.Done,
						CreatedAt:   now,
						DeadLine:    deadLine,
					},
				},
			}

			err := mongoRepo.UpdateTasks(ctx, updateReqBody)
			require.NoError(t, err)
		})
		t.Run("FetchAllTasksDetails", func(t *testing.T) {
			result, err := mongoRepo.FetchAllTasksDetails(ctx, now.Add(-time.Second*2), now)
			require.NoError(t, err)
			require.Len(t, result, 2)
		})
		t.Run("DeleteTaskFromID", func(t *testing.T) {
			err := mongoRepo.DeleteTaskFromID(ctx, task1ID.Hex())
			require.NoError(t, err)

			result, err := mongoRepo.FetchTasks(ctx, "", 0)
			require.NoError(t, err)
			require.Len(t, result, 2)
		})
	})

	t.Run("Test User Functions", func(t *testing.T) {
		testUserEmail := "g@gmail.com"
		testUserName := "Gaurav"
		testUser := &model.User{
			Name:           testUserName,
			Password:       "gaurav12345",
			Email:          testUserEmail,
			TaskCollection: testCollection,
		}
		t.Run("CreateUser", func(t *testing.T) {
			err := mongoRepo.CreateUser(ctx, testUser)
			require.NoError(t, err)
		})

		t.Run("FetchUserByEmail", func(t *testing.T) {
			user, err := mongoRepo.FetchUserByEmail(ctx, testUserEmail)
			require.NoError(t, err)
			require.NotNil(t, user)
			require.Equal(t, testUserName, user.Name)
		})

		t.Run("FetchAllUsersData", func(t *testing.T) {
			res, err := mongoRepo.FetchAllUsersData(ctx)
			require.NoError(t, err)
			require.Len(t, res, 1)
		})
	})
}
