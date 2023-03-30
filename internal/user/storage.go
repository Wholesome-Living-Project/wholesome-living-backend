package user

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// how the user is stored in the database
type userDB struct {
	FirstName   string `json:"firstName" bson:"firstName"`
	LastName    string `json:"lastName" bson:"lastName"`
	DateOfBirth string `json:"dateOfBirth" bson:"dateOfBirth"`
	Email       string `json:"email" bson:"email"`
	CreatedAt   string `json:"createdAt" bson:"createdAt"`
	ID          string `json:"id" bson:"id"`
}

type Storage struct {
	db *mongo.Database
}

func NewStorage(db *mongo.Database) *Storage {
	return &Storage{
		db: db,
	}
}

func (s *Storage) create(createUserObject createUserRequest, ctx context.Context) (string, error) {
	collection := s.db.Collection("users")

	createdAt := time.Now().Format("2006-01-02 15:04:05")

	insertObj := userDB{
		FirstName:   createUserObject.FirstName,
		LastName:    createUserObject.LastName,
		DateOfBirth: createUserObject.DateOfBirth,
		Email:       createUserObject.Email,
		CreatedAt:   createdAt,
		ID:          createUserObject.ID,
	}

	result, err := collection.InsertOne(ctx, insertObj)
	if err != nil {
		return "", err
	}

	// convert the object id to a string
	return result.InsertedID.(primitive.ObjectID).Hex(), nil
}

func (s *Storage) get(id string, ctx context.Context) (userDB, error) {
	collection := s.db.Collection("users")
	result, err := collection.Find(ctx, bson.M{"id": id})
	user := userDB{}
	if err != nil {
		return user, err
	}

	if err = result.All(ctx, &user); err != nil {
		return user, err
	}

	return user, nil

}

func (s *Storage) getAll(ctx context.Context) ([]userDB, error) {
	collection := s.db.Collection("users")

	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}

	users := make([]userDB, 0)
	if err = cursor.All(ctx, &users); err != nil {
		return nil, err
	}

	return users, nil
}
