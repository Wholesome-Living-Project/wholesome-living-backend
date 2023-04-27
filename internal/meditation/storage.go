package meditation

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

type MeditationDB struct {
	UserID         string `json:"userId" bson:"userId"`
	MeditationTime int    `json:"meditationTime" bson:"meditationTime"`
	EndTime        int64  `json:"endTime" bson:"endTime"`
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

	createdAt := time.Now().Unix()

	insertObj := MeditationDB{
		UserID:         request.UserID,
		MeditationTime: request.MeditationTime,
		EndTime:        createdAt,
	}

	result, err := collection.InsertOne(ctx, insertObj)

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
		fmt.Println("error in find")
		return nil, err
	}

	meditations := make([]MeditationDB, 0)
	for cursor.Next(ctx) {
		var meditation MeditationDB
		if err := cursor.Decode(&meditation); err != nil {
			fmt.Println("error in decode")
			return nil, err
		}
		meditations = append(meditations, meditation)
	}

	// return the meditation list
	return meditations, nil
}
