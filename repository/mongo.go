package repository

import (
	"context"
	"errors"
	"task-manager/constants"
	"task-manager/dto"
	"task-manager/model"
	"task-manager/myerrors"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoRepository struct {
	Client     *mongo.Client
	Collection string
}

// function to create tasks
func (m *MongoRepository) CreateTasks(ctx context.Context, tasks *[]model.Task) error {

	documents := make([]interface{}, len(*tasks))
	for i, task := range *tasks {
		documents[i] = task
	}

	_, err := m.Client.Database(constants.TaskManagerDB).Collection(m.Collection).InsertMany(ctx, documents)
	if err != nil {
		return err
	}
	return nil
}

// function to get tasks
func (m *MongoRepository) FetchTasks(ctx context.Context, status string, priority int) ([]*model.Task, error) {
	var tasks []*model.Task

	filter := bson.D{}

	if status != "" {
		additionalFilter := bson.D{
			{
				Key:   constants.Status,
				Value: status,
			},
		}

		filter = append(filter, additionalFilter...)
	}

	// -1 -> descending, 1 -> ascending
	findOptions := options.Find()
	findOptions.SetSort(bson.D{{Key: constants.Priority, Value: -1}}) // default desc
	if priority != 0 {
		findOptions.SetSort(bson.D{{Key: constants.Priority, Value: priority}})
	}

	cursor, err := m.Client.Database(constants.TaskManagerDB).Collection(m.Collection).Find(ctx, filter, findOptions)
	if err != nil {
		return nil, err
	}

	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var task *model.Task
		if err := cursor.Decode(&task); err != nil {
			return nil, err
		}
		tasks = append(tasks, task)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return tasks, nil
}

// function to delete task from id
func (m *MongoRepository) DeleteTaskFromID(ctx context.Context, id string) error {

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	filter := bson.D{
		{
			Key:   constants.MongoID,
			Value: objectID,
		},
	}

	result, err := m.Client.Database(constants.TaskManagerDB).Collection(m.Collection).DeleteOne(ctx, filter)
	if err != nil {
		return err
	}

	if result.DeletedCount == 0 {
		return errors.New("no document found to delete with the id passed")
	}
	return nil
}

func (m *MongoRepository) UpdateTasks(ctx context.Context, tasks *dto.TaskCreateRequestBody) error {

	for _, task := range *tasks.Tasks {
		filter := bson.D{
			{
				Key:   constants.MongoID,
				Value: task.ID,
			},
		}
		update := bson.D{
			{
				Key: "$set",
				Value: bson.D{
					{Key: "title", Value: task.Title},
					{Key: "description", Value: task.Description},
					{Key: "priority", Value: task.Priority},
					{Key: "status", Value: task.Status},
					{Key: "created_at", Value: task.CreatedAt},
					{Key: "dead_line", Value: task.DeadLine},
				},
			},
		}
		if _, err := m.Client.Database(constants.TaskManagerDB).Collection(m.Collection).UpdateOne(ctx, filter, update); err != nil {
			return err
		}
	}

	return nil
}

// function to the existence of the user by email
func (m *MongoRepository) FetchUserByEmail(ctx context.Context, email string) (*model.User, error) {

	filter := bson.D{
		{
			Key:   "email",
			Value: email,
		},
	}

	var u *model.User
	res := m.Client.Database(constants.UserDB).Collection(constants.UserCollection).FindOne(ctx, filter)

	if err := res.Decode(&u); err != nil && err != mongo.ErrNoDocuments {
		return nil, err
	}

	if res.Err() == mongo.ErrNoDocuments {
		return nil, errors.New(myerrors.NoDocumentsErr)
	}
	return u, nil
}

// function to create the user
func (m *MongoRepository) CreateUser(ctx context.Context, u *model.User) error {
	_, err := m.Client.Database(constants.UserDB).Collection(constants.UserCollection).InsertOne(ctx, u)
	return err
}

// function to create the collection
func (m *MongoRepository) CreateTaskCollection(ctx context.Context, collectionName string) error {
	return m.Client.Database(constants.TaskManagerDB).CreateCollection(ctx, collectionName)
}
