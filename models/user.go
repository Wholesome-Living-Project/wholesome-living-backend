package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type User struct {
	ID          primitive.ObjectID `json:"id" bson:"_id"`
	FirstName   string             `json:"title" bson:"title"`
	LastName    string             `json:"completed" bson:"completed"`
	CreatedAt   string             `json:"description" bson:"description"`
	DateOfBirth string             `json:"date" bson:"date"`
}
