package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Task struct {
	ID           primitive.ObjectID `json:"id" bson:"_id"`
	Title        string             `json:"title" bson:"title"`
	Descritption string             `json:"descritption" bson:"descritption"`
	Priority     int                `json:"priority" bson:"priority"`
	Status       string             `json:"status" bson:"status"`
	CreatedAt    time.Time          `json:"created_at" bson:"created_at"`
	DeadLine     time.Time          `json:"dead_line" bson:"dead_line"`
}

type TaskCreateRequestBody struct {
	Tasks []Task `json:"tasks"`
}
