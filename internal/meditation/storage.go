package meditation

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type MeditationRecord struct {
	ID             string `json:"id" bson:"_id"`
	UserID         string `json:"userId" bson:"userId"`
	MeditationTime string `json:"meditationTime" bson:"meditationTime"`
	EndTime        string `json:"endTime" bson:"endTime"`
}

// how a meditation is stored in the database
type meditationDB struct {
	ID   primitive.ObjectID `bson:"_id" json:"id"`
	Name string             `bson:"name" json:"name"`
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
	collection := s.db.Collection("mediation")

	result, err := collection.InsertOne(ctx, request)
	if err != nil {
		return "", err
	}

	// convert the object id to a string
	return result.InsertedID.(primitive.ObjectID).Hex(), nil
}

func (s *Storage) get(meditationID string, ctx context.Context) (MeditationRecord, error) {
	collection := s.db.Collection("mediation")
	meditationRecord := MeditationRecord{}

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

func (s *Storage) getAll(ctx context.Context) ([]meditationDB, error) {
	collection := s.db.Collection("meditation")

	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}

	meditations := make([]meditationDB, 0)
	if err = cursor.All(ctx, &meditations); err != nil {
		return nil, err
	}

	return meditations, nil
}
