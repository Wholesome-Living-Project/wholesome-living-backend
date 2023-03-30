package meditation

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// how the user is stored in the database
type meditationDB struct {
	ID   primitive.ObjectID `bson:"_id" json:"id"`
	Name string             `bson:"name" json:"name"`
}

type MediationStorage struct {
	db *mongo.Database
}

func NewMediationStorage(db *mongo.Database) *MediationStorage {
	return &MediationStorage{
		db: db,
	}
}

func (s *MediationStorage) createMediation(name string, ctx context.Context) (string, error) {
	collection := s.db.Collection("mediation")

	result, err := collection.InsertOne(ctx, bson.M{"name": name})
	if err != nil {
		return "", err
	}

	// convert the object id to a string
	return result.InsertedID.(primitive.ObjectID).Hex(), nil
}

func (s *MediationStorage) getAllMeditations(ctx context.Context) ([]meditationDB, error) {
	collection := s.db.Collection("meditation")

	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}

	users := make([]meditationDB, 0)
	if err = cursor.All(ctx, &meditations); err != nil {
		return nil, err
	}

	return meditations, nil
}