package meditation

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type meditationDB struct {
	ID             string `json:"id" bson:"_id"`
	UserID         string `json:"userId" bson:"userId"`
	MeditationTime string `json:"meditationTime" bson:"meditationTime"`
	EndTime        string `json:"endTime" bson:"endTime"`
}

type Storage struct {
	db *mongo.Database
}

func NewStorage(db *mongo.Database) *Storage {
	return &Storage{
		db: db,
	}
}

func (s *Storage) create(request createMeditationRequest, ctx context.Context) (string, error) {
	collection := s.db.Collection("meditation")
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

func (s *Storage) get(meditationID string, ctx context.Context) (meditationDB, error) {
	collection := s.db.Collection("meditation")
	meditationRecord := meditationDB{}

	objectID, err := primitive.ObjectIDFromHex(meditationID)
	if err != nil {
		return meditationRecord, err
	}

	cursor := collection.FindOne(ctx, bson.M{"_id": objectID})

	if err := cursor.Err(); err != nil {
		return meditationRecord, err
	}

	if err := cursor.Decode(&meditationRecord); err != nil {
		return meditationRecord, err
	}
	return meditationRecord, nil
}

func (s *Storage) getAllOfOneUser(userID string, ctx context.Context) ([]meditationDB, error) {
	collection := s.db.Collection("meditation")
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

	meditations := make([]meditationDB, 0)
	for cursor.Next(ctx) {
		var meditation meditationDB
		if err := cursor.Decode(&meditation); err != nil {
			return nil, err
		}
		meditations = append(meditations, meditation)
	}
	if err := cursor.Err(); err != nil {
		return nil, err
	}

	// return the meditation list
	return meditations, nil
}
