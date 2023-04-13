package finance

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type spendingDB struct {
	ID           string `json:"id" bson:"_id"`
	UserID       string `json:"userId" bson:"userId"`
	SpendingTime string `json:"spendingTime" bson:"spendingTime"`
	EndTime      string `json:"endTime" bson:"endTime"`
}

type Storage struct {
	db *mongo.Database
}

func NewStorage(db *mongo.Database) *Storage {
	return &Storage{
		db: db,
	}
}

func (s *Storage) create(request createSpendingRequest, ctx context.Context) (string, error) {
	collection := s.db.Collection("Spending")
	userCollection := s.db.Collection("users")

	//Check if user exists
	userResult := userCollection.FindOne(ctx, bson.M{"_id": request.UserID})

	if err := userResult.Err(); err != nil {
		return "", err
	}

	result, err := collection.InsertOne(ctx, request)

	if err != nil {
		return "", err
	}

	// convert the object id to a string
	return result.InsertedID.(primitive.ObjectID).Hex(), nil
}

func (s *Storage) get(SpendingID string, ctx context.Context) (spendingDB, error) {
	collection := s.db.Collection("Spending")
	db := spendingDB{}

	objectID, err := primitive.ObjectIDFromHex(SpendingID)
	if err != nil {
		return db, err
	}

	cursor := collection.FindOne(ctx, bson.M{"_id": objectID})

	if err := cursor.Err(); err != nil {
		return db, err
	}

	if err := cursor.Decode(&db); err != nil {
		return db, err
	}
	return db, nil
}

func (s *Storage) getAllOfOneUser(userID string, ctx context.Context) ([]spendingDB, error) {
	collection := s.db.Collection("Spending")
	userCollection := s.db.Collection("users")

	//Check if user exists
	userResult := userCollection.FindOne(ctx, bson.M{"_id": userID})

	if err := userResult.Err(); err != nil {
		fmt.Println("Error finding user:", err)
		return nil, err
	}

	cursor, err := collection.Find(ctx, bson.M{"userId": userID})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	spendings := make([]spendingDB, 0)
	for cursor.Next(ctx) {
		var Spending spendingDB
		if err := cursor.Decode(&Spending); err != nil {
			return nil, err
		}
		spendings = append(spendings, Spending)
	}
	if err := cursor.Err(); err != nil {
		return nil, err
	}

	// return the Spending list
	return spendings, nil
}
