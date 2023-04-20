package meditation

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type MeditationDB struct {
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

func (s *Storage) Create(request createMeditationRequest, ctx context.Context) (string, error) {
	collection := s.db.Collection("meditation")

	result, err := collection.InsertOne(ctx, request)

	if err != nil {
		return "", err
	}

	// convert the object id to a string
	return result.InsertedID.(primitive.ObjectID).Hex(), nil
}

func (s *Storage) Get(meditationID string, ctx context.Context) (MeditationDB, error) {
	collection := s.db.Collection("meditation")
	meditationRecord := MeditationDB{}

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

func (s *Storage) GetAllOfOneUser(userID string, ctx context.Context) ([]MeditationDB, error) {
	collection := s.db.Collection("meditation")

	cursor, err := collection.Find(ctx, bson.M{"userId": userID})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	meditations := make([]MeditationDB, 0)
	for cursor.Next(ctx) {
		var meditation MeditationDB
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
