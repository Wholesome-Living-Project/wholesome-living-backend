package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Meditation struct {
	ID             primitive.ObjectID `json:"id" bson:"_id"`
	UserID         primitive.ObjectID `json:"user_id" bson:"user_id"`
	LastName       string             `json:"completed" bson:"completed"`
	CreatedAt      string             `json:"date" bson:"date"`
	MeditationTime string             `json:"meditation_time" bson:"meditation_time"`
}
