package model

import "go.mongodb.org/mongo-driver/bson/primitive"

type Task struct {
	ID          primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Title       string             `json:"title" bson:"title"`
	Description string             `json:"description" bson:"description"`
	Priority    int                `json:"priority" bson:"priority"`
	Status      string             `json:"status" bson:"status"`
	CreatedAt   string             `json:"created_at" bson:"created_at"`
	DeadLine    string             `json:"dead_line" bson:"dead_line"`
}
