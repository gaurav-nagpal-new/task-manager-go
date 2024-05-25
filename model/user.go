package model

import "go.mongodb.org/mongo-driver/bson/primitive"

type User struct {
	ID             primitive.ObjectID `bson:"_id,omitempty" json:"-"`
	Name           string             `bson:"name" json:"name"`
	Password       string             `bson:"password" json:"password"`
	Email          string             `bson:"email" json:"email"`
	TaskCollection string             `bson:"task_collection"`
}
