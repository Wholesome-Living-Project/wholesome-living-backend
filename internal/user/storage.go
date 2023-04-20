package user

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// how the user is stored in the database
type UserDB struct {
	FirstName   string       `json:"firstName" bson:"firstName"`
	LastName    string       `json:"lastName" bson:"lastName"`
	DateOfBirth string       `json:"dateOfBirth" bson:"dateOfBirth"`
	Email       string       `json:"email" bson:"email"`
	CreatedAt   string       `json:"createdAt" bson:"createdAt"`
	ID          string       `json:"id" bson:"_id"`
	Plugins     []pluginType `json:"plugins" bson:"plugins"`
}

type Storage struct {
	db *mongo.Database
}

func NewStorage(db *mongo.Database) *Storage {
	return &Storage{
		db: db,
	}
}

func (s *Storage) Create(createUserObject createUserRequest, ctx context.Context) (string, error) {
	collection := s.db.Collection("users")

	createdAt := time.Now().Format("2006-01-02 15:04:05")
	var plugins []pluginType

	insertObj := UserDB{
		FirstName:   createUserObject.FirstName,
		LastName:    createUserObject.LastName,
		DateOfBirth: createUserObject.DateOfBirth,
		Email:       createUserObject.Email,
		CreatedAt:   createdAt,
		ID:          createUserObject.ID,
		Plugins:     plugins,
	}

	result, err := collection.InsertOne(ctx, insertObj)
	if err != nil {
		fmt.Println(err)
		return "", err
	}

	// convert the object id to a string
	return result.InsertedID.(string), nil
}

func (s *Storage) Get(id string, ctx context.Context) (UserDB, error) {
	collection := s.db.Collection("users")
	result := collection.FindOne(ctx, bson.M{"_id": id})
	user := UserDB{}

	if result.Err() != nil {
		return user, result.Err()
	}

	if err := result.Decode(&user); err != nil {
		return user, err
	}

	return user, nil

}

func (s *Storage) GetAll(ctx context.Context) ([]UserDB, error) {
	collection := s.db.Collection("users")

	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}

	users := make([]UserDB, 0)
	if err = cursor.All(ctx, &users); err != nil {
		return nil, err
	}

	return users, nil
}

func (s *Storage) Update(user UserDB, ctx context.Context) (UserDB, error) {
	collection := s.db.Collection("users")
	result := collection.FindOneAndUpdate(ctx, bson.M{"_id": user.ID}, bson.M{"$set": bson.M{"firstName": user.FirstName, "lastName": user.LastName, "dateOfBirth": user.DateOfBirth, "email": user.Email, "plugins": user.Plugins}}, nil)

	if result.Err() != nil {
		return user, result.Err()
	}

	if err := result.Decode(&user); err != nil {
		return user, err
	}

	return user, nil

}
