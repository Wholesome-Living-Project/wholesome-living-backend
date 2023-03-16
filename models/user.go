package models

type User struct {
	ID          string `json:"id" bson:"_id"`
	FirstName   string `json:"title" bson:"title"`
	LastName    string `json:"completed" bson:"completed"`
	CreatedAt   string `json:"description" bson:"description"`
	DateOfBirth string `json:"date" bson:"date"`
}
